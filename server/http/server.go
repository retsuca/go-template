package server

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"go-template/internal/config"
	"go-template/pkg/metrics"
	"go-template/pkg/tracer"
	"go-template/server/http/handler"
	"go-template/server/http/middleware"
	"go-template/server/http/routes"

	httpclient "go-template/internal/clients/httpClient"
	logger "go-template/pkg/logger"

	http_metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	http_metrics_middleware "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

type Server struct {
	Metrics *handler.Metrics
}

// CreateHTTPServer initializes and starts an HTTP server with the given configuration.
// It sets up middleware, routes, and handles graceful shutdown.
func CreateHTPPServer(ctx context.Context, host, port string, gwMux *runtime.ServeMux) {
	// Initialize tracer
	tp, err := tracer.NewTracer()
	if err != nil {
		logger.Fatal("Failed to initialize tracer", zap.Error(err))
	}

	defer func() {
		if err := tp.Shutdown(ctx); err != nil {
			logger.Error("Error shutting down tracer provider", zap.Error(err))
		}
	}()

	r := chi.NewRouter()

	// Setup middleware
	setupMiddleware(r)

	// Create handler instance
	h := createHandler()

	// Setup routes
	routes.SetupRoutes(r, h, gwMux)

	// Start server
	startServer(ctx, r, host, port)
}

func createHandler() *handler.Handler {
	client := httpclient.NewClient(httpclient.ClientOptions{
		BaseURL: &url.URL{
			Scheme: "https",
			Host:   "hacker-news.firebaseio.com",
			Path:   "v0/",
		},
		InsecureSkipVerify: false,
	})

	metrics := registerMetrics()

	return handler.NewHandler(client, metrics)
}

func registerMetrics() *handler.Metrics {
	metrics := &handler.Metrics{
		HelloCounter: metrics.NewCounterVec("hello_counter_http", []string{"hello"}, ""),
		HelloGauge:   metrics.NewGaugeVec("hello_gauge_http", []string{"hello"}, ""),
	}

	return metrics
}

// setupMiddleware configures all middleware for the server.
func setupMiddleware(r *chi.Mux) {
	appName := config.Get(config.APP_NAME)

	// Add tracing middleware
	r.Use(func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, appName, otelhttp.WithFilter(otelReqFilter))
	})

	// Add error handling middleware
	r.Use(middleware.ErrorHandler)

	// Add recover middleware
	r.Use(middleware.Recoverer)

	// Add request ID middleware
	r.Use(chimiddleware.RequestID)

	// Add logger middleware
	r.Use(chimiddleware.Logger)

	// Add CORS middleware
	r.Use(middleware.DefaultCORS().Handler)

	// Add prometheus middleware
	mdlw := http_metrics_middleware.New(http_metrics_middleware.Config{
		Recorder: http_metrics.NewRecorder(http_metrics.Config{}),
	})
	r.Use(std.HandlerProvider("", mdlw))
}

func otelReqFilter(req *http.Request) bool {
	return req.URL.Path != "/metrics"
}

// startServer starts the HTTP server and handles graceful shutdown.
func startServer(ctx context.Context, r *chi.Mux, host, port string) {
	// Create server context with cancellation
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	server := &http.Server{
		Addr:    host + ":" + port,
		Handler: r,
	}

	// Start server in a goroutine
	go func() {
		logger.Info("Starting HTTP server", zap.String("address", server.Addr))
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Fatal error in HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	logger.Info("Shutting down HTTP server...")

	// Give the server 10 seconds to complete pending requests
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during server shutdown", zap.Error(err))
	}

	logger.Info("HTTP server shutdown complete")
}
