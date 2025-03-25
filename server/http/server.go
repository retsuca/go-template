package server

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo-contrib/echoprometheus"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware" // echo-swagger middleware
	echoSwagger "github.com/swaggo/echo-swagger"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "go-template/docs"
	httpclient "go-template/internal/clients/httpClient"
	"go-template/internal/config"
	logger "go-template/pkg/logger"
	"go-template/pkg/metrics"
	"go-template/pkg/tracer"
	swagger "go-template/proto/gen/swagger"
	"go-template/server/http/handler"
	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
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

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true

	// Setup middleware
	setupMiddleware(e)

	// Create handler instance
	h := createHandler()

	// Setup routes with handler instance
	setupRoutes(e, gwMux)

	// Setup handlers
	setupHandlers(e, h)

	// Start server
	startServer(ctx, e, host, port)
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
func setupMiddleware(e *echo.Echo) {
	appName := config.Get(config.APP_NAME)
	metricName := strings.ReplaceAll(appName, "-", "_")

	// Add tracing middleware first to ensure all requests are traced
	e.Use(otelecho.Middleware(appName,
		otelecho.WithSkipper(func(c echo.Context) bool {
			return c.Path() == "/metrics" || c.Path() == "/health"
		}),
	))

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(echoprometheus.NewMiddleware(metricName))

	// Add request ID middleware
	e.Use(middleware.RequestID())
}

// setupRoutes configures all routes for the server.
func setupRoutes(e *echo.Echo, gwMux *runtime.ServeMux) {
	e.GET("/metrics", echoprometheus.NewHandler())

	// Configure Swagger UI only when gRPC gateway is enabled
	if gwMux != nil {
		// Serve gRPC-Gateway Swagger documentation
		e.GET("/swagger/swagger.json", func(c echo.Context) error {
			return c.JSONBlob(http.StatusOK, swagger.ApidocsSwaggerJson)
		})

		e.GET("/swagger/*", echo.WrapHandler(httpSwagger.Handler(
			httpSwagger.URL("/swagger/swagger.json"), // The url pointing to API definition
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		)))
		e.Any("/*", echo.WrapHandler(gwMux))
	} else {
		e.GET("/swagger/*", echoSwagger.WrapHandler)
	}

	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})
}

// setupHandlers initializes and configures all handlers.
func setupHandlers(e *echo.Echo, h *handler.Handler) {
	e.GET("/", h.Hello)
	e.GET("/withparam", h.HelloWithParam)
}

// startServer starts the HTTP server and handles graceful shutdown.
func startServer(ctx context.Context, e *echo.Echo, host, port string) {
	// Create server context with cancellation
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// Start server in a goroutine
	go func() {
		addr := host + ":" + port
		logger.Info("Starting HTTP server", zap.String("address", addr))

		if err := e.Start(addr); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal("Fatal error in HTTP server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal
	<-ctx.Done()
	logger.Info("Shutting down HTTP server...")

	// Give the server 10 seconds to complete pending requests
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		logger.Error("Error during server shutdown", zap.Error(err))
	}

	logger.Info("HTTP server shutdown complete")
}
