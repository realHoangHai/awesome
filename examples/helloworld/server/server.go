// Package main provide an implementation example of gRPC using awesome.
package main

import (
	"context"
	"github.com/realHoangHai/awesome/config"
	pb "github.com/realHoangHai/awesome/examples/helloworld/helloworld"
	"github.com/realHoangHai/awesome/internal/server"
	"github.com/realHoangHai/awesome/pkg/log"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type (
	service struct {
		pb.UnimplementedGreeterServer
	}
)

// SayHello implements pb.GreeterServer interface.
func (s *service) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Context(ctx).Info("name", req.Name)
	//if req.Name == "" {
	//	return nil, status.InvalidArgument("name must not be empty")
	//}
	return &pb.HelloReply{
		Message: "Hello " + req.GetName(),
	}, nil
}

// Register implements server.Service interface.
func (s *service) Register(srv *grpc.Server) {
	pb.RegisterGreeterServer(srv, s)
}

// RegisterWithEndpoint implements server.EndpointService interface.
func (s *service) RegisterWithEndpoint(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) {
	pb.RegisterGreeterHandlerFromEndpoint(ctx, mux, addr, opts)
}

func main() {
	err := log.Init(log.FromEnv(config.WithFileNoError(".env")))
	if err != nil {
		return
	}

	srv := &service{}
	if err := server.New(
		server.FromEnv(config.WithFileNoError(".env")),
	).ListenAndServe(srv); err != nil {
		log.Panic(err)
	}
}
