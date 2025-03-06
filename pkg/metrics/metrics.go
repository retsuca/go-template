package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	opsProcessed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "myapp_processed_ops_total",
		// Subsystem: "Subsystem",
		// Namespace: "Namespace",
		Help: "The total number of processed events",
	})

	gauge = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: "golang",
			Name:      "my_gauge",
			Help:      "This is my gauge",
		})

	histogram = promauto.NewHistogram(
		prometheus.HistogramOpts{
			Namespace: "golang",
			Name:      "my_histogram",
			Help:      "This is my histogram",
		})

	summary = promauto.NewSummary(
		prometheus.SummaryOpts{
			Namespace: "golang",
			Name:      "my_summary",
			Help:      "This is my summary",
		})
)

// func init() {
// 	opsProcessed.Inc()
// 	gauge.Set(1)
// 	histogram.Observe(1)
// 	summary.Observe(1)
// 	summary.Observe(1)
// 	summary.Observe(1)
// 	summary.Observe(1)
// 	summary.Observe(1)
// 	summary.Observe(1)
// 	summary.Observe(4)
// }
