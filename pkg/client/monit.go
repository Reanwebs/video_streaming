package client

import (
	"context"
	"fmt"
	clientinterfaces "videoStreaming/pkg/client/clientInterfaces"
	"videoStreaming/pkg/config"
	"videoStreaming/pkg/domain"
	"videoStreaming/pkg/pb/monit"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type monitClient struct {
	Server monit.MonitizationServiceClient
}

func InitClient(c *config.Config) (clientinterfaces.MonitClient, error) {
	cc, err := grpc.Dial(c.Product_SVC, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return NewAuthClient(monit.NewMonitizationServiceClient(cc)), nil
}

func NewAuthClient(server monit.MonitizationServiceClient) clientinterfaces.MonitClient {
	return &monitClient{
		Server: server,
	}
}

func (m *monitClient) HealthCheck(ctx context.Context) (string, error) {
	return "", nil
}

func (m *monitClient) VideoReward(ctx context.Context, request domain.VideoRewardRequest) error {
	fmt.Println("VideoReward")
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
