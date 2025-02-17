package controllers

import (
	"go-template/pkg/tracer"
	"net/http"

	"github.com/labstack/echo/v4"
)

func Hello(c echo.Context) error {
	tracer.TestTrace(c.Request().Context())
	return c.String(http.StatusOK, "Hello, World!")
}
