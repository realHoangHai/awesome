# awesome

Actually, it is not awesome at all. It's name is awesome because it's owner 
using JetBrains Goland and he was too lazy to change the default idea project
name.

## Getting Started

### Usage

#### [Grpc Gateway](https://github.com/grpc-ecosystem/grpc-gateway)

The following instructions assume you are using Go Modules for dependency 
management. Use a [tool](./pkg/tools/tools.go) dependency to track the versions of the following 
executable packages:

	// +build tools
	
	package tools
	
	import (
		_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
		_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
		_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
		_ "google.golang.org/protobuf/cmd/protoc-gen-go"
	)

Run go mod tidy to resolve the versions. Install by running
	
	make init

#### [Ent](https://entgo.io/) - an easy-to-use ORM which developed by Facebook

For each entity you want to expose, `cd` to [repo](./internal/storage) and create a
new file in the [ent/schema](internal/storage/ent/schema) directory by

 	ent init --target ./ent/schema/ <EntityName>

Configing `Fields` and `Edges` in file already generated.

Finally, run `go generate ./ent` to generate the code.

<details>

<summary>Note</summary>

- [generate.go](./internal/repo/ent/generate.go) is a helper to generate entities.
  It is already exists before.

- You are still in the [repo](./internal/storage) directory while `init` and `gen`. If not, you must change
  the command

</details>


### Start your own

Create new gRPC service

```go
func (s *service) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloReply, error) {
    return &pb.HelloReply{
        Message: "Hello " + req.GetName(),
    }, nil
}

// Register implements server.Service interface
// It registers gRPC APIs with gRPC server.
func (s *service) Register(srv *grpc.Server) {
    pb.RegisterGreeterServer(srv, s)
}

// RegisterWithEndpoint implements server.EndpointService interface
// It is used to expose REST API using gRPC Gateway.
func (s *service) RegisterWithEndpoint(ctx context.Context, mux *runtime.ServeMux, addr string, opts []grpc.DialOption) {
    pb.RegisterGreeterHandlerFromEndpoint(ctx, mux, addr, opts)
}
```

Start a simple server, get configurations from environment variables.

```go
package main

import (
    "github.com/realHoangHai/awesome/server"
)

func main() {
    srv := &service{}
    if err := server.ListenAndServe(srv); err != nil {
        panic(err)
    }
}
```

More complex with custom options.

```go
package main

import (
    "github.com/realHoangHai/awesome/pkg/log"
    "github.com/realHoangHai/awesome/server"
)

func main() {
    srv := server.New(
        server.FromEnv(),
        server.PProf(""),
        server.Address(":8088"),
        server.JWT("secret"),
        server.Web("/", "web", "index.html"),
        server.Logger(log.Fields("service", "my_service")),
        server.CORS(true, []string{"*"}, []string{"POST"}, []string{"http://localhost:8088"}),
    )
    if err := srv.ListenAndServe( /*services...*/ ); err != nil {
        panic(err)
    }
}

```

## Features

Currently, awesomeProject supports following features:

### Server

- Exposes both gRPC and REST in 1 single port.
- Authentication interceptors
- Other options: CORS, HTTP Handler, Interceptors,...


### Cache

- Standard cache service interface.
- Memory cache.
- Redis plugin.

### Config

- Standard config interface.
- Config from environment variables.
- Config from file and other options.

### Log

- Standard logger interface.
- Logrus implementation.
- Interceptors for HTTP & gRPC.

### Util

- Some utilities that might need during the development using awesome.

Interceptors: [go-grpc-middleware](https://github.com/grpc-ecosystem/go-grpc-middleware)