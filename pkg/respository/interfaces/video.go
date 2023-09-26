package interfaces

import "videoStreaming/pkg/domain"

type VideoRepo interface {
	CreateVideoid(domain.ToSaveVideo) (string, error)
}
