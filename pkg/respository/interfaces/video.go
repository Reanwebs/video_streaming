package interfaces

import "videoStreaming/pkg/pb"

type VideoRepo interface {
	CreateVideoid(string, string) (string, error)
	FindAllVideo() ([]*pb.VideoID, error)
}
