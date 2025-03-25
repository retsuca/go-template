package handler

import (
	"github.com/prometheus/client_golang/prometheus"
	pbName "go-template/proto/gen/go/helloservice/v1/name"
)

type Metrics struct {
	HelloCounter *prometheus.CounterVec
	HelloGauge   *prometheus.GaugeVec
}

type HelloServer struct {
	pbName.UnimplementedGreeterServiceServer
	Metrics *Metrics
}

func NewHelloServer(metrics *Metrics) *HelloServer {
	return &HelloServer{
		Metrics: metrics,
	}
}
