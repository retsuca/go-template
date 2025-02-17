package server

import (
	"errors"
	logger "go-template/pkg/logger"
	"go-template/server/controllers"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"github.com/labstack/echo-contrib/echoprometheus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateHTPPServer(host, port string) {

	e := echo.New()

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
