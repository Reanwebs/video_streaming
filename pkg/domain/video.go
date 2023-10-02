package domain

import "time"

type Video struct {
	ID             uint `gorm:"primarykey"`
	Video_id       string
	Archived       bool `gorm:"default:false"`
	S3_path        string
	User_name      string
	Avatar_id      string
	Title          string
	Discription    string
	Interest       string
	Thumbnail_id   string
	Views          uint
	Starred        uint
	UserId         string
	Exclusive      bool
	Coin_for_watch uint `gorm:"default:0"`
	Blocked        bool `gorm:"default:false"`
}

type ToSaveVideo struct {
	S3Path      string `json:"s3path"`
	UserName    string `json:"userName"`
	AvatarId    string `json:"avatarId"`
	Title       string `json:"title"`
	Discription string `json:"discription"`
	Intrest     string `json:"interest"`
	ThumbnailId string `json:"thumbnailId"`
	Video_id    string
	UserId      string
}

type Star struct {
	ID       uint `gorm:"primarykey"`
	VideoID  string
	UserName string
}

type Viewer struct {
	ID        uint `gorm:"primarykey"`
	VideoID   string
	UserName  string
	Timestamp time.Time
}

type BlockedVideo struct {
	ID        uint `gorm:"primarykey"`
	VideoID   string
	Reason    string
	Timestamp time.Time
}

type ReportVideo struct {
	ID      uint `gorm:"primarykey"`
	VideoId string
	Reason  string
}

type ReportedVideo struct {
	Reason         string
	Video_id       string
	Archived       bool `gorm:"default:false"`
	S3_path        string
	User_name      string
	Avatar_id      string
	Title          string
	Discription    string
	Interest       string
	Thumbnail_id   string
	Views          uint
	Starred        uint
	UserId         string
	Exclusive      bool
	Coin_for_watch uint `gorm:"default:0"`
	Blocked        bool `gorm:"default:false"`
}
