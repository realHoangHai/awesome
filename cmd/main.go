package main

import (
	"github.com/realHoangHai/awesome/config"
	"github.com/realHoangHai/awesome/internal/server"
	"github.com/realHoangHai/awesome/pkg/log"
)

func main() {
	cfg, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	var services []server.Service
	userService, closedb, err := wireApp(&cfg)
	defer closedb()
	if err != nil {
		log.Fatal(err)
	}

	services = append(services, userService)

	s := server.New(server.FromEnv(&cfg))
	if err := s.Run(services...); err != nil {
		log.Fatal(err)
	}
}
