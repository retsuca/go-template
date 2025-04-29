package handler

import (
	"context"

	"go-template/pkg/logger"

	"go.uber.org/zap"

	pbName "go-template/proto/gen/go/helloservice/v1/name"
)

func (s *HelloServer) SayHello(_ context.Context, in *pbName.SayHelloRequest) (*pbName.SayHelloResponse, error) {
	logger.Info("Received: ", zap.String("name", in.GetName()))

	s.Metrics.HelloCounter.WithLabelValues("test").Inc()

	return &pbName.SayHelloResponse{Message: "Hello " + in.GetName()}, nil
}
