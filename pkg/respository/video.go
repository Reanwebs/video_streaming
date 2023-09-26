package repository

import (
	"videoStreaming/pkg/domain"
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

func (c *videoRepo) CreateVideoid(input domain.ToSaveVideo) (string, error) {
	// Create a new Video record using the input data
	video := &domain.Video{
		S3Path:      input.S3Path,
		UserName:    input.UserName,
		AvatarId:    input.AvatarId,
		Title:       input.Title,
		Discription: input.Discription,
		Interest:    input.Interest,
		ThumbnailId: input.ThumbnailId,
	}

	if err := c.DB.Create(video).Error; err != nil {
		return "", err
	}

	videoId := video.VideoId

	return videoId, nil
}

// func (c *videoRepo) FindAllVideo() ([]*pb.VideoID, error) {
// 	var videoid []*pb.VideoID
// 	if err := c.DB.Model(&domain.Video{}).Find(&videoid).Error; err != nil {
// 		return nil, err
// 	}
// 	return videoid, nil
// }
