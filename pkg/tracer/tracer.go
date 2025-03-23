// Package tracer provides OpenTelemetry tracing functionality for the application
package tracer

import (
	"context"
	"fmt"
	"time"

	"go-template/internal/config"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TracerProvider holds the tracer provider instance and shutdown function
type TracerProvider struct {
	provider *sdktrace.TracerProvider
	tracer   oteltrace.Tracer
}

// NewTracer initializes a new tracer provider with the given service name
func NewTracer() (*TracerProvider, error) {
	serviceName := config.Get(config.APP_NAME)

	// Create resource with service information
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
			semconv.DeploymentEnvironmentKey.String(config.Get(config.ENV)),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// Create span processor based on environment
	bsp, err := createBatchSpanProcessor()
	if err != nil {
		return nil, fmt.Errorf("failed to create span processor: %w", err)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// Set global tracer provider and propagator
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Create tracer instance
	tracer := tp.Tracer(serviceName)

	return &TracerProvider{
		provider: tp,
		tracer:   tracer,
	}, nil
}

// Shutdown gracefully shuts down the tracer provider
func (tp *TracerProvider) Shutdown(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := tp.provider.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown tracer provider: %w", err)
	}

	return nil
}

// Tracer returns the tracer instance
func (tp *TracerProvider) Tracer() oteltrace.Tracer {
	return tp.tracer
}

// createBatchSpanProcessor creates a span processor based on the environment
func createBatchSpanProcessor() (sdktrace.SpanProcessor, error) {
	if config.Get(config.ENV) == "local" {
		exporter, err := stdout.New(stdout.WithPrettyPrint())
		if err != nil {
			return nil, fmt.Errorf("failed to create stdout exporter: %w", err)
		}
		return sdktrace.NewBatchSpanProcessor(exporter), nil
	}

	endpoint := config.Get(config.OTEL_EXPORTER_OTLP_ENDPOINT)
	if endpoint == "" {
		return nil, fmt.Errorf("OTEL_EXPORTER_OTLP_ENDPOINT is not set")
	}

	ctx := context.Background()
	conn, err := grpc.Dial(endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithGRPCConn(conn),
		otlptracegrpc.WithTimeout(5*time.Second),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create OTLP exporter: %w", err)
	}

	return sdktrace.NewBatchSpanProcessor(exporter), nil
}

// StartSpan starts a new span with the given name and attributes
func StartSpan(ctx context.Context, name string, attrs ...attribute.KeyValue) (context.Context, oteltrace.Span) {
	return otel.Tracer("").Start(ctx, name,
		oteltrace.WithAttributes(attrs...),
	)
}

// SpanFromContext returns the current span from context
func SpanFromContext(ctx context.Context) oteltrace.Span {
	return oteltrace.SpanFromContext(ctx)
}
