package server_test

import (
	"context"
	"github.com/realHoangHai/awesome/config"
	"github.com/realHoangHai/awesome/internal/server"
	"github.com/realHoangHai/awesome/pkg/log"
)

func ExampleListenAndServe() {
	cfg, _ := config.LoadConfig("../../config/config.yaml")
	if err := server.ListenAndServe(&cfg /*services ...Service*/); err != nil {
		panic(err)
	}
}

func ExampleListenAndServeContext() {
	cfg, _ := config.LoadConfig("config.toml")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	if err := server.ListenAndServeContext(ctx, &cfg /*, services ...Service*/); err != nil {
		panic(err)
	}
}

func ExampleNew_fromEnvironmentVariables() {
	cfg, _ := config.LoadConfig("config.toml")
	srv := server.New(server.FromEnv(&cfg))
	if err := srv.ListenAndServe( /*services ...Service*/ ); err != nil {
		log.Panic(err)
	}
}

func ExampleNew_withOptions() {
	srv := server.New(
		server.Address(":8080"),
		server.Logger(log.Fields("service", "awesome")),
	)
	if err := srv.ListenAndServe( /*services ...Service*/ ); err != nil {
		panic(err)
	}
}
