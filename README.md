# awesome 🚀

Actually, it is not awesome at all. It's name is awesome because it's owner 
using JetBrains Goland and he was too lazy to change the default idea project
name.

## Getting Started

### Table of contents

- [Project Layout](#project-layout)
- [Commands](#commands)
- [Usages](#usages)
- [Features](#features)

### Project Layout
awesome uses the following project layout

```
.
├── api                       contains all .proto and compiled go file of each modules
│   └── user                  user module
│       └── v1                version of user module
├── certs                     certificate for ssl connection
├── cmd                       main applications of the project
│   ├── main.go               read from config and run the application
│   ├── wire.go               wire the application
│   └── wire_gen.go           wire generated file
├── config                    configuration files for different environments
├── internal                  private application and library code
│   ├── auth                  authentication feature
│   ├── biz                   business logic layer of the project
│   │   ├── biz.go            provider of business logic layer
│   │   ├── user.go           user business
│   │   └── ...               other entity business
│   ├── health                healthcheck feature
│   ├── server                configuration server grpc and http
│   ├── service               transport layer of the project
│   │   ├── user.go           user service 
│   │   └── ...               other entity service 
│   └── storage               storage layer of the project
│       ├── ent               ent for storage layer 
│       │   ├── schema        contains all schema of the entity and its relation
│       │   ├── generate.go   to generate the code of the entity
│       │   └── ...           genreated code of the entity 
│       ├── store             provides storage layer for different modules 
│       ├── user              user storage 
│       └── ...               other entity storage 
├── pkg                       public library code
│   ├── encoding              encoding lib
│   ├── log                   structured and context-aware logger
│   ├── jwt                   json web token
│   ├── status                wrapped status code for grpc
│   ├── tools                 structured and context-aware logger
│   └── utils                 contains some useful functions
├── third_party               third party library for protocol buffers
├── .gitignore                .gitignore
├── app.toml                  configuration file for the project
├── docker-compose.yml        docker compose file
├── Dockerfile                dockerfile for building the project
├── Makefile                  makefile for building the project
└── README.md                 readme file
```

### Commands

The following instructions assume you are using Go Modules for dependency 
management. Use a [tool](./pkg/tools/tools.go) dependency to track the versions of the following 
executable packages:

```
// +build tools

package tools

import (
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
    _ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
    _ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
    _ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
```

Run go mod tidy to resolve the versions. Install by running

```
make init
```

Custom ent orm for query and migrate the database.

```shell
# for each entity want to expose
ent init --target .internal/storage/ent/schema/ <EntityName>

# generate the ent code for the entity
make generate
```


### Usages

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
	"github.com/realHoangHai/awesome/internal/server"
)

func main() {
    srv := &service{}
    if err := server.Run(srv); err != nil {
        panic(err)
    }
}
```

More complex with custom options.

```go
package main

import (
  "github.com/realHoangHai/awesome/config"
  "github.com/realHoangHai/awesome/internal/server"
  "github.com/realHoangHai/awesome/pkg/log"
)

func main() {
  cfg, _ := config.LoadConfig(".")
  srv := server.New(
    server.FromEnv(&cfg),
    server.PProf(""),
    server.Address(":8088"),
    server.JWT("secret"),
    server.Web("/", "web", "index.html"),
    server.Logger(log.Fields("service", "my_service")),
    server.CORS(true, []string{"*"}, []string{"POST"}, []string{"http://localhost:8088"}),
  )
  if err := srv.Run( /*services...*/); err != nil {
    panic(err)
  }
}

```

## Features

Currently, awesomeProject supports following features:

### Server

- Exposes both gRPC and REST in 1 single port.
- Internal APIs:
  - [Prometheus](https://github.com/grpc-ecosystem/go-grpc-prometheus) metrics.
  - [Health](https://github.com/grpc/grpc/blob/master/doc/health-checking.md) checks.
  - Debug profiling.
- Authentication interceptors
- Other options: CORS, HTTP Handler, Serving Single Page Application, Interceptors,...

### Auth

- Authenticator interface.
- JWT
- Authenticator, WhiteList, Chains.
- Interceptors for both gRPC & HTTP

### Health

- Health check for readiness and liveness.
- Utilities for checking health.

### Config

- Standard config interface.
- Config from environment variables.
- Config from file and other options.

### Log

- Standard logger interface.
- [Logrus](https://github.com/sirupsen/logrus) implementation.
- Interceptors for HTTP & gRPC.