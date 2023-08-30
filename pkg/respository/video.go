package repository

import (
	"videoStreaming/pkg/domain"
	"videoStreaming/pkg/pb"
	"videoStreaming/pkg/respository/interfaces"

	"gorm.io/gorm"
)

type videoRepo struct {
	DB *gorm.DB
}

func NewVideoRepo(db *gorm.DB) interfaces.VideoRepo {
	return &videoRepo{
		DB: db,
	}
}

func (c *videoRepo) CreateVideoid(videoid string, path string) (string, error) {
	err := c.DB.Create(&domain.Video{
		VideoId: videoid,
		S3Path:  path,
	}).Error
	if err != nil {
		return "", err
	}
	return videoid, nil
}

func (c *videoRepo) FindAllVideo() ([]*pb.VideoID, error) {
	var videoid []*pb.VideoID
	if err := c.DB.Model(&domain.Video{}).Find(&videoid).Error; err != nil {
		return nil, err
	}
	return videoid, nil
}
