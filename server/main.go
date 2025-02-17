package server

import (
	"errors"
	logger "go-template/pkg/logger"
	"go-template/server/controllers"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"

	"github.com/labstack/echo-contrib/echoprometheus"
	stdout "go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateHTPPServer(host, port string) {

	e := echo.New()

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

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware("my-server"))

	e.Use(echoprometheus.NewMiddleware("gotemplate"))
	e.GET("/metrics", echoprometheus.NewHandler())

	e.GET("/", controllers.Hello)

	if err := e.Start(host + ":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.Fatalw("Fatal error http server ", err)
	}

}
