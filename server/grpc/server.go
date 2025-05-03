package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"go-template/pkg/logger"
	"go-template/pkg/metrics"
	"go-template/pkg/tracer"
	"go-template/server/grpc/config"
	"go-template/server/grpc/gateway"
	"go-template/server/grpc/handler"
	"go-template/server/grpc/health"
	httpServer "go-template/server/http"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"

	pbName "go-template/proto/gen/go/helloservice/v1/name"
)

// Server represents a gRPC server instance.
type Server struct {
	grpcServer   *grpc.Server
	config       *config.Config
	shutdownChan chan os.Signal
	tp           *tracer.TracerProvider
	Metrics      *handler.Metrics
	health       *health.Health
	gateway      *gateway.Gateway
}

// NewServer creates a new Server instance with the given configuration.
func NewServer(cfg *config.Config, metrics *handler.Metrics) *Server {
	if cfg == nil {
		return nil
	}

	return &Server{
		config:       cfg,
		shutdownChan: make(chan os.Signal, 1),
		Metrics:      metrics,
		health:       health.NewHealth(),
		gateway:      gateway.NewGateway(),
	}
}

// Start initializes and starts both GRPC and HTTP servers.
func (s *Server) Start(ctx context.Context) error {
	logger.Info("Initializing server",
		zap.String("host", s.config.Host),
		zap.String("grpc_port", s.config.GRPCPort),
		zap.String("http_port", s.config.HTTPPort))

	// Initialize tracer
	tp, err := tracer.NewTracer()
	if err != nil {
		return fmt.Errorf("failed to initialize tracer: %w", err)
	}
	s.tp = tp

	// Setup signal handling
	signal.Notify(s.shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Initialize gRPC server
	if err := s.initGRPCServer(); err != nil {
		return fmt.Errorf("failed to initialize gRPC server: %w", err)
	}

	// Setup gRPC-Gateway
	if err := s.gateway.Setup(ctx, s.config.Host, s.config.GRPCPort); err != nil {
		return fmt.Errorf("failed to setup gRPC gateway: %w", err)
	}

	// Start servers
	wg := &sync.WaitGroup{}
	errChan := make(chan error, 2)

	// Start gRPC server
	s.startGRPCServer(wg, errChan)

	// Start HTTP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		httpServer.CreateHTPPServer(ctx, s.config.Host, s.config.HTTPPort, s.gateway.GetMux())
	}()

	// Wait for shutdown signal or error
	select {
	case err := <-errChan:
		return fmt.Errorf("server error: %w", err)
	case sig := <-s.shutdownChan:
		logger.Info("Received shutdown signal", zap.String("signal", sig.String()))
		s.gracefulShutdown(ctx)
	case <-ctx.Done():
		logger.Info("Context canceled, initiating shutdown")
		s.gracefulShutdown(ctx)
	}

	wg.Wait()
	return nil
}

// initGRPCServer initializes the gRPC server and registers services.
func (s *Server) initGRPCServer() error {
	// Create gRPC server with tracing instrumentation
	otelHandler := otelgrpc.NewServerHandler()
	s.grpcServer = grpc.NewServer(
		grpc.StatsHandler(otelHandler),
	)

	// Register services
	helloServer := handler.NewHelloServer(s.Metrics)
	pbName.RegisterGreeterServiceServer(s.grpcServer, helloServer)

	// Register health check service
	s.health.Register(s.grpcServer)
	s.health.SetServingStatus("", healthpb.HealthCheckResponse_SERVING)

	// Register reflection service
	reflection.Register(s.grpcServer)

	addr := net.JoinHostPort(s.config.Host, s.config.GRPCPort)
	logger.Info("gRPC server initialized", zap.String("address", addr))

	return nil
}

// startGRPCServer starts the gRPC server in a goroutine.
func (s *Server) startGRPCServer(wg *sync.WaitGroup, errChan chan<- error) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		addr := net.JoinHostPort(s.config.Host, s.config.GRPCPort)
		lis, err := net.Listen("tcp", addr)
		if err != nil {
			logger.Error("Failed to listen on gRPC port", zap.Error(err), zap.String("address", addr))
			errChan <- fmt.Errorf("failed to listen on gRPC port: %w", err)
			return
		}

		logger.Info("Starting gRPC server", zap.String("address", addr))

		if err := s.grpcServer.Serve(lis); err != nil {
			logger.Error("gRPC server error", zap.Error(err))
			errChan <- fmt.Errorf("gRPC server error: %w", err)
		}
	}()
}

// gracefulShutdown handles graceful shutdown of the server.
func (s *Server) gracefulShutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Shutdown tracer provider
	if s.tp != nil {
		if err := s.tp.Shutdown(shutdownCtx); err != nil {
			logger.Error("Error shutting down tracer provider", zap.Error(err))
		}
	}

	// Shutdown gRPC server
	logger.Info("Shutting down gRPC server")

	stopped := make(chan struct{})

	go func() {
		s.grpcServer.GracefulStop()
		close(stopped)
	}()

	select {
	case <-shutdownCtx.Done():
		logger.Warn("Graceful shutdown timed out, forcing gRPC server stop")
		s.grpcServer.Stop()
	case <-stopped:
		logger.Info("gRPC server stopped gracefully")
	}

	logger.Info("Servers shutdown complete")
}

// CreateGRPCServer creates and starts the gRPC and HTTP servers with the given configuration.
func CreateGRPCServer(ctx context.Context, host, grpcPort, httpPort string) {
	cfg := config.NewConfig(host, grpcPort, httpPort)
	metrics := registerMetrics()

	server := NewServer(cfg, metrics)
	if err := server.Start(ctx); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}

func registerMetrics() *handler.Metrics {
	metrics := &handler.Metrics{
		HelloCounter: metrics.NewCounterVec("hello_counter_grpc", []string{"hello"}, ""),
		HelloGauge:   metrics.NewGaugeVec("hello_gauge_grpc", []string{"hello"}, ""),
	}

	return metrics
}
