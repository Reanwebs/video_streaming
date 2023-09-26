package repository

import (
	"errors"
	"fmt"
	"strings"
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

func (c *videoRepo) CreateVideoid(input domain.ToSaveVideo) (bool, error) {
	// Clean up the input.UserName by removing unwanted symbols
	cleanedUserName := strings.Replace(input.UserName, "\"", "", -1)

	// Create a new Video record using the cleaned username
	video := &domain.Video{
		S3_path:      input.S3Path,
		User_name:    cleanedUserName,
		Avatar_id:    input.AvatarId,
		Title:        input.Title,
		Discription:  input.Discription,
		Interest:     input.Intrest,
		Thumbnail_id: input.ThumbnailId,
	}

	if err := c.DB.Create(video).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (c *videoRepo) FetchUserVideos(userName string) ([]*domain.Video, error) {
	var data []*domain.Video
	if err := c.DB.Model(&domain.Video{}).
		Where("user_name = ? AND archived = ?", userName, false).
		Find(&data).
		Error; err != nil {
		return nil, err
	}

	if len(data) == 0 {
		fmt.Println("fetching empty array")
		return []*domain.Video{}, errors.New("there is novideo")
	}

	return data, nil
}
