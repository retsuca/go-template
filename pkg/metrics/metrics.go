package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func NewCounterVec(metricName string, labels []string, help string) *prometheus.CounterVec {
	return promauto.NewCounterVec(prometheus.CounterOpts{
		Name: metricName,
		Help: help,
	}, labels)
}

func NewGaugeVec(metricName string, labels []string, help string) *prometheus.GaugeVec {
	return promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: metricName,
		Help: help,
	}, labels)
}

func NewHistogramVec(metricName string, labels []string, help string) *prometheus.HistogramVec {
	return promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: metricName,
		Help: help,
	}, labels)
}

func NewSummaryVec(metricName string, labels []string, help string) *prometheus.SummaryVec {
	return promauto.NewSummaryVec(prometheus.SummaryOpts{
		Name: metricName,
		Help: help,
	}, labels)
}
