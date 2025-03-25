package handler

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
)

type Client interface {
	Do(ctx context.Context, method, path string, body any, args map[string]string) ([]byte, error)
}

type Metrics struct {
	HelloCounter *prometheus.CounterVec
	HelloGauge   *prometheus.GaugeVec
}
type Handler struct {
	HTTPClient Client
	Metrics    *Metrics
}

func NewHandler(client Client, metrics *Metrics) *Handler {
	return &Handler{
		HTTPClient: client,
		Metrics:    metrics,
	}
}
