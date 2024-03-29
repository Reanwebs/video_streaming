package client

import (
	"context"
	"log"
	clientinterfaces "videoStreaming/pkg/client/clientInterfaces"
	"videoStreaming/pkg/config"
	"videoStreaming/pkg/domain"
	"videoStreaming/pkg/pb/monit"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type monitClient struct {
	Server monit.MonitizationClient
}

func InitClient(c *config.Config) (clientinterfaces.MonitClient, error) {
	cc, err := grpc.Dial(c.MONIT_SVC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return NewMonitClient(monit.NewMonitizationClient(cc)), nil
}

func NewMonitClient(server monit.MonitizationClient) clientinterfaces.MonitClient {
	return &monitClient{
		Server: server,
	}
}

func (m *monitClient) HealthCheck(ctx context.Context) (string, error) {
	return "", nil
}

func (m *monitClient) VideoReward(ctx context.Context, request domain.VideoRewardRequest) error {
	_, err := m.Server.VideoReward(ctx, &monit.VideoRewardRequest{
		UserID:    request.UserID,
		VideoID:   request.VideoID,
		Reason:    "view",
		Views:     request.Views,
		PaidCoins: 0,
	})
	if err != nil {
		return err
	}
	return nil
}

func (m *monitClient) ExclusiveContent(ctx context.Context, request domain.ExclusiveContentRequest) error {
	_, err := m.Server.ExclusiveContent(ctx, &monit.ExclusiveContentRequest{
		UserID:    request.UserID,
		VideoID:   request.VideoID,
		Reason:    "paid",
		Owner:     request.Owner,
		PaidCoins: request.PaidCoins,
	})
	if err != nil {
		return err
	}
	return nil
}
