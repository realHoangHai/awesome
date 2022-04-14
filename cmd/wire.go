//go:build wireinject
// +build wireinject

// The build tag makes sure the stub is not built in the final build.
package main

import (
	"context"
	"github.com/google/wire"
	"github.com/realHoangHai/awesome/config"
	"github.com/realHoangHai/awesome/internal/biz"
	"github.com/realHoangHai/awesome/internal/repo"
	"github.com/realHoangHai/awesome/internal/service"
)

// wireApp init application.
func InitializeServer(ctx context.Context, config *config.Config) {
	wire.Build(
		repo.ProviderSet,
		biz.ProviderSet,
		service.ProviderSet,
	)
}
