package main

import (
	"context"
	pb "github.com/realHoangHai/awesomeProject/examples/helloworld/helloworld"
	"github.com/realHoangHai/awesomeProject/pkg/log"
	"github.com/realHoangHai/awesomeProject/server"
	"github.com/realHoangHai/awesomeProject/status"

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
	if req.Name == "" {
		return nil, status.InvalidArgument("name must not be empty")
	}
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
	srv := &service{}
	opts := []server.Option{
		server.FromEnv(),
		server.Web("/", "public", "index.html"),
	}
	if err := server.New(opts...).ListenAndServe(srv); err != nil {
		log.Panic(err)
	}
}
