package interfaces

import (
	"videoStreaming/pkg/domain"
	"videoStreaming/pkg/pb"
)

type VideoRepo interface {
	CreateVideoid(domain.ToSaveVideo) (bool, error)
	FetchUserVideos(string) ([]*domain.Video, error)
	FindArchivedVideos(userName string) ([]*domain.Video, error)
	ArchivedVideos(VideoId uint) (bool, error)
	FetchAllVideos() ([]*domain.Video, error)
	GetVideoById(string, string) (*domain.Video, bool, []*domain.Video, error)
	ToggleStar(string, string, bool) (bool, error)
	BlockVideo(domain.BlockedVideo) (bool, error)
	GetReportedVideos() ([]domain.ReportedVideo, error)
	ReportVideo(input *pb.ReportVideoRequest) (bool, error)
	FetchExclusiveVideos() ([]*domain.Video, error)
}
