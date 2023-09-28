package repository

import (
	"fmt"
	"time"
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

	video := &domain.Video{
		S3_path:      input.S3Path,
		User_name:    input.UserName,
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

	if err := c.DB.
		Where("user_name = ?", userName).
		Find(&data).
		Error; err != nil {
		return nil, err
	}

	if len(data) == 0 {
		fmt.Println("Fetching empty array")
		return []*domain.Video{}, nil
	}
	return data, nil
}

func (c *videoRepo) FindArchivedVideos(userName string) ([]*domain.Video, error) {
	var data []*domain.Video
	if err := c.DB.Model(&domain.Video{}).
		Where("user_name = ? AND archived = ?", userName, true).
		Find(&data).
		Error; err != nil {
		return nil, err
	}

	if len(data) == 0 {
		fmt.Println("fetching empty array")
		return []*domain.Video{}, nil
	}

	return data, nil
}

func (c *videoRepo) ArchivedVideos(VideoId uint) (bool, error) {
	var video domain.Video
	if err := c.DB.First(&video, VideoId).Error; err != nil {
		return false, err
	}

	video.Archived = !video.Archived

	if err := c.DB.Save(&video).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (c *videoRepo) FetchAllVideos() ([]*domain.Video, error) {
	var data []*domain.Video
	if err := c.DB.Model(&domain.Video{}).
		Where("archived = ?", false).
		Find(&data).
		Error; err != nil {
		return nil, err
	}

	if len(data) == 0 {
		fmt.Println("Fetching empty array")
		return []*domain.Video{}, nil
	}

	return data, nil
}

func (c *videoRepo) GetVideoById(id uint, userName string) (*domain.Video, bool, error) {

	var video domain.Video

	if err := c.DB.Where("id = ?", id).First(&video).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, fmt.Errorf("Video with ID %d not found", id)
		}
		return nil, false, err
	}

	var star domain.Star
	if err := c.DB.Where("video_id = ? AND user_name = ?", id, userName).First(&star).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return &video, false, nil
		}
		return nil, false, err
	}

	video.Views++
	if err := c.DB.Save(&video).Error; err != nil {
		return nil, false, err
	}

	view := domain.Viewer{
		VideoID:   id,
		UserName:  userName,
		Timestamp: time.Now(),
	}
	if err := c.DB.Create(&view).Error; err != nil {
		return nil, false, err
	}

	return &video, true, nil
}

func (c *videoRepo) ToggleStar(id uint, userName string, starred bool) (bool, error) {
	if starred {
		star := domain.Star{
			VideoID:  id,
			UserName: userName,
		}
		if err := c.DB.Create(&star).Error; err != nil {
			return false, err
		}
	} else {
		if err := c.DB.Where("video_id = ? AND user_name = ?", id, userName).Delete(&domain.Star{}).Error; err != nil {
			return false, err
		}
	}
	return true, nil
}
