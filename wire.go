// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/steebchen/keskin-api/api"
	"github.com/steebchen/keskin-api/handlers"
	"github.com/steebchen/keskin-api/prisma"
	"github.com/steebchen/keskin-api/server"
)

func Initialize() (*server.Server, error) {
	wire.Build(
		prisma.ProviderSet,
		api.ProviderSet,
		server.ProviderSet,
		handlers.ProviderSet,
	)
	return &server.Server{}, nil
}
