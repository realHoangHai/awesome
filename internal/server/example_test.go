package server_test

import (
	"github.com/realHoangHai/awesome/config"
	"github.com/realHoangHai/awesome/internal/server"
	"github.com/realHoangHai/awesome/pkg/log"
)

func ExampleNew() {
	cfg, _ := config.LoadConfig("config.toml")
	srv := server.New(server.FromEnv(&cfg))
	if err := srv.Run( /*services ...Service*/ ); err != nil {
		log.Panic(err)
	}
}

func ExampleNew_withOptions() {
	srv := server.New(
		server.Address(":8080"),
		server.Logger(log.Fields("service", "awesome")),
	)
	if err := srv.Run( /*services ...Service*/ ); err != nil {
		panic(err)
	}
}
