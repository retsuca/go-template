package tracer

import (
	"context"
	"go-template/pkg/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var otelSDKTrace *sdktrace.TracerProvider
var tracer = otel.Tracer("echo-server")

func init() {
	exporter, err := stdout.New(stdout.WithPrettyPrint())
	if err != nil {
		logger.Fatalw("Fatal error initing Tracer ", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithBatcher(exporter),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

	otelSDKTrace = tp
}

func TestTrace(ctx context.Context) {
	_, span := tracer.Start(ctx, "getUser", oteltrace.WithAttributes(attribute.String("id", "222")))
	defer span.End()

}
