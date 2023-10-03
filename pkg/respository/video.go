package repository

import (
	"errors"
	"fmt"
	"log"
	"time"
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

func (c *videoRepo) CreateVideoid(input domain.ToSaveVideo) (bool, error) {

	video := &domain.Video{
		S3_path:        input.S3Path,
		User_name:      input.UserName,
		Avatar_id:      input.AvatarId,
		Title:          input.Title,
		Discription:    input.Discription,
		Interest:       input.Intrest,
		Thumbnail_id:   input.ThumbnailId,
		Views:          0,
		Starred:        0,
		Video_id:       input.Video_id,
		UserId:         input.UserId,
		Exclusive:      input.Exclusive,
		Coin_for_watch: input.Coin_for_watch,
	}

	if err := c.DB.Create(video).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (c *videoRepo) FetchUserVideos(userName string) ([]*domain.Video, error) {
	var data []*domain.Video

	if err := c.DB.
		Where("user_name = ? AND blocked = ?", userName, false).
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
		Where("archived = ? AND blocked = ?", false, false).
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

func (c *videoRepo) GetVideoById(id string, userName string) (*domain.Video, bool, error) {
	// Find the video by its VideoId
	var video domain.Video
	if err := c.DB.Where("video_id = ?", id).First(&video).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, fmt.Errorf("Video with ID %s not found", id)
		}
		return nil, false, err
	}

	var star domain.Star
	if err := c.DB.Where("video_id = ? AND user_name = ?", id, userName).First(&star).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// User has not starred this video
			return &video, false, nil
		}
		return nil, false, err
	}

	// Increment the view count
	video.Views++
	if err := c.DB.Save(&video).Error; err != nil {
		fmt.Println("error in Increment the view count")
		return nil, false, err
	}
	fmt.Println(" Increment the view count")

	// Record the viewer's information
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

func (c *videoRepo) ToggleStar(id string, userName string, starred bool) (bool, error) {
	// Start a transaction to ensure data consistency
	tx := c.DB.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create or delete a star record
	if starred {
		star := domain.Star{
			VideoID:  id,
			UserName: userName,
		}
		if err := tx.Create(&star).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	} else {
		if err := tx.Where("video_id = ? AND user_name = ?", id, userName).Delete(&domain.Star{}).Error; err != nil {
			tx.Rollback()
			return false, err
		}
	}

	// Update the video's Starred field
	var video domain.Video
	if err := tx.Where("video_id = ?", id).First(&video).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if starred {
		video.Starred++
	} else {
		if video.Starred > 0 {
			video.Starred--
		}
	}

	if err := tx.Save(&video).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	// Commit the transaction if everything succeeded
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return false, err
	}

	return true, nil
}

func (c *videoRepo) BlockVideo(input domain.BlockedVideo) (bool, error) {
	tx := c.DB.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}

	var video domain.Video
	if err := tx.Where("video_id = ? AND blocked = true", input.VideoID).First(&video).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return false, err
		}
	} else {
		log.Printf("Video with ID %s is already blocked\n", input.VideoID)
		tx.Rollback()
		return false, nil
	}

	if err := tx.Model(&domain.Video{}).Where("video_id = ?", input.VideoID).Update("blocked", false).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Where("video_id = ?", input.VideoID).Delete(&domain.BlockedVideo{}).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	blockRecord := &domain.BlockedVideo{
		VideoID:   input.VideoID,
		Reason:    input.Reason,
		Timestamp: time.Now(),
	}

	if err := tx.Create(blockRecord).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	if err := tx.Commit().Error; err != nil {
		return false, err
	}

	return true, nil
}

func (c *videoRepo) ReportVideo(input *pb.ReportVideoRequest) (bool, error) {
	data := &domain.ReportVideo{
		VideoId: input.VideoId,
		Reason:  input.Reason,
	}

	if err := c.DB.Create(data).Error; err != nil {
		return false, err
	}

	return true, nil
}

func (c *videoRepo) GetReportedVideos() ([]domain.ReportedVideo, error) {
	var reportedVideos []domain.ReportedVideo

	if err := c.DB.Table("videos").
		Joins("JOIN report_videos ON videos.Video_id = report_videos.VideoId").
		Scan(&reportedVideos).Error; err != nil {
		return nil, err
	}

	return reportedVideos, nil
}
