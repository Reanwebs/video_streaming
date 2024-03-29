//go:build wireinject
// +build wireinject

package di

import (
	api "videoStreaming/pkg/api"
	"videoStreaming/pkg/api/service"
	"videoStreaming/pkg/client"
	"videoStreaming/pkg/config"
	"videoStreaming/pkg/db"
	repository "videoStreaming/pkg/respository"

	"github.com/google/wire"
)

func InitializeServe(c *config.Config) (*api.Server, error) {
	wire.Build(
		db.Initdb,
		repository.NewVideoRepo,
		client.InitClient,
		service.NewVideoServer,
		api.NewgrpcServe,
	)
	return &api.Server{}, nil
}
