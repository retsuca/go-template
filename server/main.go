package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go-template/internal/config"
	"go-template/pkg/log"
	"go-template/server/controllers"
	"go.uber.org/zap"
)

func init() {
	e := echo.New()

	// Middleware
	e.Use(log.ZapLogger(zap.L()))
	e.Use(middleware.Recover())

	e.GET("/test", controllers.GetTest)

	e.Logger.Fatal(e.Start(":" + config.Port))

}
