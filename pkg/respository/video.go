package repository

import (
	"errors"
	"fmt"
	"strings"
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
		Where("archived = ? AND blocked = ? AND exclusive = ?", false, false, false).
		Find(&data).
		Error; err != nil {
		return nil, err
	}
	fmt.Println("\nall")
	if len(data) == 0 {
		fmt.Println("Fetching empty array")
		return []*domain.Video{}, nil
	}

	return data, nil
}

func (c *videoRepo) GetVideoById(id string, userName string) (*domain.Video, bool, []*domain.Video, error) {
	var video domain.Video
	isStarred := true
	if err := c.DB.Where("video_id = ?", id).First(&video).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, false, nil, fmt.Errorf("Video with ID %s not found", id)
		}
		return nil, false, nil, err
	}

	var star domain.Star
	err := c.DB.Where("video_id = ? AND user_name = ?", id, userName).First(&star).Error
	if err == gorm.ErrRecordNotFound || err != nil {
		isStarred = false
	}

	video.Views++
	if err := c.DB.Save(&video).Error; err != nil {
		fmt.Println("error in Increment the view count")
		return nil, false, nil, err
	}
	fmt.Println(" Increment the view count")

	view := domain.Viewer{
		VideoID:   id,
		UserName:  userName,
		Timestamp: time.Now(),
	}
	if err := c.DB.Create(&view).Error; err != nil {
		return nil, false, nil, err
	}

	var suggestions []*domain.Video
	if err := c.DB.Model(&domain.Video{}).
		Where("archived = ? AND blocked = ? AND exclusive = ? AND interest = ?", false, false, false, video.Interest).
		Find(&suggestions).
		Error; err != nil {
		return nil, false, nil, err
	}

	if len(suggestions) > 3 {
		suggestions = suggestions[:3]
	}

	return &video, isStarred, suggestions, nil
}

func (c *videoRepo) ToggleStar(id string, userName string, starred bool) (bool, error) {
	tx := c.DB.Begin()
	if tx.Error != nil {
		return false, tx.Error
	}
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	var video domain.Video
	if err := tx.Where("video_id = ?", id).First(&video).Error; err != nil {
		tx.Rollback()
		return false, err
	}

	var star domain.Star
	err := c.DB.Where("video_id = ? AND user_name = ?", id, userName).First(&star).Error
	if err != nil && !strings.Contains(err.Error(), "record not found") {
		tx.Rollback()
		return false, err
	}

	if starred && err == nil {
	} else if starred && err != nil {
		// Video is not starred, so star it
		newStar := domain.Star{
			VideoID:  id,
			UserName: userName,
		}
		if err := tx.Create(&newStar).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		video.Starred++
	} else if !starred && err == nil {
		if err := tx.Delete(&star).Error; err != nil {
			tx.Rollback()
			return false, err
		}
		if video.Starred > 0 {
			video.Starred--
		}
	} else {
		// This case should not occur
		tx.Rollback()
		return false, errors.ErrUnsupported
	}

	if err := tx.Save(&video).Error; err != nil {
		tx.Rollback()
		return false, err
	}

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
	if err := tx.Where("video_id = ?", input.VideoID).First(&video).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			tx.Rollback()
			return false, err
		}
	} else {
		if !video.Blocked {
			if err := tx.Model(&domain.Video{}).Where("video_id = ?", input.VideoID).Update("blocked", true).Error; err != nil {
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
			if err := tx.Where("video_id = ?", input.VideoID).Delete(&domain.ReportVideo{}).Error; err != nil {
				tx.Rollback()
				return false, err
			}
		} else {
			if err := tx.Model(&domain.Video{}).Where("video_id = ?", input.VideoID).Update("blocked", false).Error; err != nil {
				tx.Rollback()
				return false, err
			}
			if err := tx.Where("video_id = ?", input.VideoID).Delete(&domain.BlockedVideo{}).Error; err != nil {
				tx.Rollback()
				return false, err
			}
		}
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
		Joins("JOIN report_videos ON videos.Video_id = report_videos.Video_id").
		Scan(&reportedVideos).Error; err != nil {
		return nil, err
	}

	return reportedVideos, nil
}

func (c *videoRepo) FetchExclusiveVideos() ([]*domain.Video, error) {
	var videos []*domain.Video
	if err := c.DB.Where("exclusive = ?", true).
		Order("coin_for_watch desc").
		Find(&videos).Error; err != nil {
		return nil, err
	}

	if len(videos) == 0 {
		fmt.Println("Fetching empty array")
		return []*domain.Video{}, nil
	}

	return videos, nil
}
