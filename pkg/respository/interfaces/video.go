package interfaces

import (
	"videoStreaming/pkg/domain"
)

type VideoRepo interface {
	CreateVideoid(domain.ToSaveVideo) (bool, error)
	FetchUserVideos(string) ([]*domain.Video, error)
	FindArchivedVideos(userName string) ([]*domain.Video, error)
	ArchivedVideos(VideoId uint) (bool, error)
	FetchAllVideos() ([]*domain.Video, error)
	GetVideoById(id uint) (*domain.Video, error)
}
