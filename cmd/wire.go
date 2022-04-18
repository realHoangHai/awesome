//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.

package main

import (
	"github.com/google/wire"
	"github.com/realHoangHai/awesome/config"
	"github.com/realHoangHai/awesome/internal/biz"
	"github.com/realHoangHai/awesome/internal/service"
	repo "github.com/realHoangHai/awesome/internal/storage"
)

// wireApp init awesome application
func wireApp(cfg *config.Config) (*service.UserService, func(), error) {
	panic(wire.Build(repo.ProviderSet, biz.ProviderSet, service.NewUserService))
}
