package gateway

import (
	"context"
	"fmt"
	"net"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pbName "go-template/proto/gen/go/helloservice/v1/name"
)

// Gateway represents the gRPC-Gateway server.
type Gateway struct {
	mux *runtime.ServeMux
}

// NewGateway creates a new Gateway instance.
func NewGateway() *Gateway {
	return &Gateway{
		mux: runtime.NewServeMux(),
	}
}

// Setup initializes and configures the gRPC-Gateway.
func (g *Gateway) Setup(ctx context.Context, host, grpcPort string) error {
	// Use new client instrumentation with propagation
	otelHandler := otelgrpc.NewClientHandler(
		otelgrpc.WithTracerProvider(otel.GetTracerProvider()),
		otelgrpc.WithPropagators(otel.GetTextMapPropagator()),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStatsHandler(otelHandler),
	}
	grpcServerEndpoint := net.JoinHostPort(host, grpcPort)

	// Register gRPC-Gateway handlers
	if err := pbName.RegisterGreeterServiceHandlerFromEndpoint(ctx, g.mux, grpcServerEndpoint, opts); err != nil {
		return fmt.Errorf("failed to register gRPC gateway: %w", err)
	}

	return nil
}

// GetMux returns the ServeMux instance.
func (g *Gateway) GetMux() *runtime.ServeMux {
	return g.mux
} 