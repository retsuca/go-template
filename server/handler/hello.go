package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-template/pkg/metrics"
	"go-template/pkg/tracer"
)

func (h *Handler) Hello(c echo.Context) error {
	metrics.OpsProcessed.Inc()
	tracer.TestTrace(c.Request().Context())

	return c.String(http.StatusOK, "Hello, World!")
}
