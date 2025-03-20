package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"go-template/pkg/metrics"
)

// @Summary      hello world
// @Description  shows hello world
// @Tags         accounts
// @Router       / [get]
func (h *Handler) Hello(c echo.Context) error {
	metrics.OpsProcessed.Inc()
	// tracer.TestTrace(c.Request().Context())

	return c.String(http.StatusOK, "Hello, World!")
}

// @Summary      hello world
// @Description  shows hello world
// @Tags         accounts
// @Router       /withparam [get]
// @Param name query string true "name"
func (h *Handler) HelloWithParam(c echo.Context) error {
	name := c.QueryParam("name")

	metrics.OpsProcessed.Inc()
	// TestTrace(c.Request().Context())

	return c.String(http.StatusOK, name)
}
