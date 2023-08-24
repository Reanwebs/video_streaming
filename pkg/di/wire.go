//go:build wireinject
// +build wireinject

package di

import (
	service "videoStreaming/pkg/api"
	api "videoStreaming/pkg/api/service"
	"videoStreaming/pkg/config"
	"videoStreaming/pkg/db"
	repository "videoStreaming/pkg/respository"

	"github.com/google/wire"
)

func InitializeServe(c *config.Config) (*api.Server, error) {
	wire.Build(db.Initdb, repository.NewVideoRepo, service.NewVideoServer, api.NewgrpcServe)
	return &api.Server{}, nil
}
