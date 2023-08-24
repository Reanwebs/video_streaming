package interfaces

import "videoStreaming/pkg/pb"

type VideoRepo interface {
	CreateVideoid(string) error
	FindAllVideo() ([]*pb.VideoID, error)
}
