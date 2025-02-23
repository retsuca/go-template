package server

import (
	"errors"
	"go-template/internal/config"
	logger "go-template/pkg/logger"
	"go-template/server/controllers"
	"net/http"
	"strings"

	"go.opentelemetry.io/contrib/instrumentation/github.com/labstack/echo/otelecho"

	"github.com/labstack/echo-contrib/echoprometheus"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func CreateHTPPServer(host, port string) {

	e := echo.New()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(otelecho.Middleware(""))

	appName := config.Get(config.APP_NAME)

	e.Use(echoprometheus.NewMiddleware(strings.Replace(appName, "-", "_", -1)))

	e.GET("/metrics", echoprometheus.NewHandler())

	e.GET("/", controllers.Hello)

	if err := e.Start(host + ":" + port); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logger.FatalErr("Fatal error http server ", err)
	}

}
