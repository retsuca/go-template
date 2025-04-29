package middleware

import (
	"net/http"

	http_metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	http_metrics_middleware "github.com/slok/go-http-metrics/middleware"
	"github.com/slok/go-http-metrics/middleware/std"
)

var mdlw = http_metrics_middleware.New(http_metrics_middleware.Config{
	Recorder: http_metrics.NewRecorder(http_metrics.Config{}),
})

func WrapMetricHandler(metricPath string, h func(w http.ResponseWriter, r *http.Request)) http.Handler {
	return std.Handler(metricPath, mdlw, http.HandlerFunc(h))
}
