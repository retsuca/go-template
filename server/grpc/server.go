package server

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	_ "go-template/docs"
	"go-template/internal/config"
	"go-template/pkg/logger"

	pbName "go-template/proto/gen/go/helloservice/v1/name"
	"go-template/server/grpc/handler"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer   *grpc.Server
	echoServer   *echo.Echo
	config       *Config
	shutdownChan chan os.Signal
}

type Config struct {
	Host     string
	GRPCPort string
	HTTPPort string
}

// NewServer creates a new Server instance with the given configuration
func NewServer(cfg *Config) *Server {
	return &Server{
		config:       cfg,
		shutdownChan: make(chan os.Signal, 1),
	}
}

// Start initializes and starts both GRPC and HTTP servers
func (s *Server) Start(ctx context.Context) error {
	// Setup signal handling
	signal.Notify(s.shutdownChan, os.Interrupt, syscall.SIGTERM)

	// Initialize servers
	s.initGRPCServer()
	s.initEchoServer()

	// Start servers
	wg := &sync.WaitGroup{}
	errChan := make(chan error, 2)

	s.startGRPCServer(wg, errChan)
	s.startEchoServer(wg, errChan)

	// Wait for shutdown signal or error
	select {
	case err := <-errChan:
		logger.Error("Server error", zap.Error(err))
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

func (s *Server) initGRPCServer() {
	s.grpcServer = grpc.NewServer()
	reflection.Register(s.grpcServer)
	pbName.RegisterGreeterServiceServer(s.grpcServer, &handler.HelloServer{})

	addr := net.JoinHostPort(s.config.Host, s.config.GRPCPort)
	logger.Info("gRPC server initialized", zap.String("address", addr))
}

func (s *Server) initEchoServer() {
	e := echo.New()

	// Add middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware(""))

	appName := config.Get(config.APP_NAME)
	e.Use(echoprometheus.NewMiddleware(strings.ReplaceAll(appName, "-", "_")))
	e.GET("/metrics", echoprometheus.NewHandler())

	// Setup gRPC-Gateway
	grpcServerEndpoint := net.JoinHostPort(s.config.Host, s.config.GRPCPort)
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}

	// Register gRPC-Gateway handlers
	gwMux := runtime.NewServeMux()
	if err := pbName.RegisterGreeterServiceHandlerFromEndpoint(context.Background(), gwMux, grpcServerEndpoint, opts); err != nil {
		logger.Error("Failed to register gRPC gateway", zap.Error(err))
	}

	// Mount gRPC-Gateway
	e.Any("/v1/*", echo.WrapHandler(gwMux))

	// Swagger routes
	e.GET("/swagger/*", echoSwagger.WrapHandler)

	s.echoServer = e
	logger.Info("HTTP server initialized", zap.String("address", net.JoinHostPort(s.config.Host, s.config.HTTPPort)))
}

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

func (s *Server) startEchoServer(wg *sync.WaitGroup, errChan chan<- error) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		addr := net.JoinHostPort(s.config.Host, s.config.HTTPPort)
		logger.Info("Starting HTTP server", zap.String("address", addr))
		if err := s.echoServer.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error("HTTP server error", zap.Error(err))
			errChan <- fmt.Errorf("HTTP server error: %w", err)
		}
	}()
}

func (s *Server) gracefulShutdown(ctx context.Context) {
	shutdownCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	// Shutdown HTTP server
	logger.Info("Shutting down HTTP server")
	if err := s.echoServer.Shutdown(shutdownCtx); err != nil {
		logger.Error("HTTP server shutdown error", zap.Error(err))
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

// CreateGRPCServer creates and starts the gRPC and HTTP servers
func CreateGRPCServer(ctx context.Context, host, grpcPort, httpPort string) {
	cfg := &Config{
		Host:     host,
		GRPCPort: grpcPort,
		HTTPPort: httpPort,
	}

	server := NewServer(cfg)
	if err := server.Start(ctx); err != nil {
		logger.Fatal("Server error", zap.Error(err))
	}
}
