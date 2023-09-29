package clientinterfaces

import (
	"context"
	"videoStreaming/pkg/domain"
)

type MonitClient interface {
	HealthCheck(context.Context) (string, error)
	VideoReward(context.Context, domain.VideoRewardRequest) error
}
