package handler

import (
	"context"
	"log"

	pbName "go-template/proto/gen/go/helloservice/v1/name"
)

type HelloServer struct {
	pbName.UnimplementedGreeterServiceServer
}

func (s *HelloServer) SayHello(_ context.Context, in *pbName.SayHelloRequest) (*pbName.SayHelloResponse, error) {
	log.Printf("Received: %v", in.GetName())
	return &pbName.SayHelloResponse{Message: "Hello " + in.GetName()}, nil
}
