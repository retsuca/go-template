package tracer

import (
	"context"
	"fmt"

	"go-template/internal/config"
	"go-template/pkg/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	oteltrace "go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var tracer = otel.Tracer("echo-server")

func init() {
	bsp := createBatchSpanProcessor()

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(config.Get(config.APP_NAME)),
	)

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
}

func TestTrace(ctx context.Context) {
	_, span := tracer.Start(ctx, "getUser", oteltrace.WithAttributes(attribute.String("id", "222")))
	defer span.End()
}

func createBatchSpanProcessor() sdktrace.SpanProcessor {
	if config.Get(config.ENV) == "local" {
		exporter, err := stdout.New(stdout.WithPrettyPrint())
		if err != nil {
			logger.FatalErr("Fatal error initing Tracer ", err)
		}

		bsp := sdktrace.NewBatchSpanProcessor(exporter)

		return bsp
	}

	conn, err := initConn()
	if err != nil {
		logger.FatalErr("Fatal error initing Tracer ", err)
	}

	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		logger.FatalErr("Fatal error initing Tracer ", err)
	}

	bsp := sdktrace.NewBatchSpanProcessor(exporter)

	return bsp
}

func initConn() (*grpc.ClientConn, error) {
	conn, err := grpc.NewClient(config.Get(config.OTEL_EXPORTER_OTLP_ENDPOINT),
		// Note the use of insecure transport here. TLS is recommended in production.
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create gRPC connection to collector: %w", err)
	}

	return conn, err
}
