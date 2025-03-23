package handler

import (
	"context"

	"go-template/pkg/logger"
	pbName "go-template/proto/gen/go/helloservice/v1/name"

	"go.uber.org/zap"
)

type HelloServer struct {
	pbName.UnimplementedGreeterServiceServer
}

func (s *HelloServer) SayHello(_ context.Context, in *pbName.SayHelloRequest) (*pbName.SayHelloResponse, error) {
	logger.Info("Received: ", zap.String("name", in.GetName()))

	return &pbName.SayHelloResponse{Message: "Hello " + in.GetName()}, nil
}
