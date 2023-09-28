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
	GetVideoById(uint, string) (*domain.Video, bool, error)
	ToggleStar(uint, string, bool) (bool, error)
}
