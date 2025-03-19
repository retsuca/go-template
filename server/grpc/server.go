package server

import (
	"context"
	"net"
	"net/http"
	"sync"

	_ "go-template/docs"
	"go-template/pkg/logger"

	pbName "go-template/proto/gen/go/helloservice/v1/name"

	"go-template/server/grpc/handler"

	v2 "github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	httpSwagger "github.com/swaggo/http-swagger"

	gorilla "github.com/gorilla/mux"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

type fnRestRegister func(ctx context.Context, mux *v2.ServeMux, endpoint string, opts []grpc.DialOption) error

func CreateGRPCServer(ctx context.Context, host, grpcPort, httpPort string) {
	lis, err := net.Listen("tcp", host+":"+grpcPort)
	if err != nil {
		logger.FatalErr("Fatal error grpc server ", err)
	}

	s := grpc.NewServer()
	reflection.Register(s)

	grpcServerEndpoint := host + ":" + grpcPort

	restRegister := []fnRestRegister{
		pbName.RegisterGreeterServiceHandlerFromEndpoint,
	}

	mux := createGeneralMux(ctx, grpcServerEndpoint, restRegister)

	pbName.RegisterGreeterServiceServer(s, &handler.HelloServer{})

	wg := sync.WaitGroup{}

	grpcErrs := make(chan error, 1)
	wg.Add(1)
	go func() {
		err := s.Serve(lis)
		grpcErrs <- err
		wg.Done()
	}()

	httpErrs := make(chan error, 1)
	wg.Add(1)
	go func() {
		http.ListenAndServe(host+":"+httpPort, mux)
		httpErrs <- err
		wg.Done()
	}()

	select {
	case err := <-grpcErrs:
		logger.FatalErr(err.Error(), err)
	case err := <-httpErrs:
		logger.FatalErr(err.Error(), err)
	case <-ctx.Done():
		logger.FatalErr(err.Error(), err)
	}

	wg.Wait()
}

func createGeneralMux(ctx context.Context, grpcServerEndpoint string, restRegister []fnRestRegister) *gorilla.Router {
	gwMux := v2.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	for _, register := range restRegister {
		register(context.Background(), gwMux, grpcServerEndpoint, opts)
	}

	mux := gorilla.NewRouter()
	mux.HandleFunc("/swagger/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./proto/gen/swagger/apidocs.swagger.json")
	})
	mux.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/swagger.json"), // The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)
	mux.PathPrefix("/").Handler(gwMux)

	return mux
}
