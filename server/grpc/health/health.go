package health

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

// Health represents the health check server.
type Health struct {
	server *health.Server
}

// NewHealth creates a new Health instance.
func NewHealth() *Health {
	return &Health{
		server: health.NewServer(),
	}
}

// Register registers the health check service with the gRPC server.
func (h *Health) Register(grpcServer *grpc.Server) {
	healthpb.RegisterHealthServer(grpcServer, h.server)
}

// SetServingStatus sets the serving status of the service.
func (h *Health) SetServingStatus(service string, status healthpb.HealthCheckResponse_ServingStatus) {
	h.server.SetServingStatus(service, status)
} 