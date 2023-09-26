package domain

type Video struct {
	ID           uint `gorm:"primarykey"`
	Archived     bool `gorm:"default:false"`
	S3_path      string
	User_name    string
	Avatar_id    string
	Title        string
	Discription  string
	Interest     string
	Thumbnail_id string
}

type ToSaveVideo struct {
	S3Path      string `json:"s3path"`
	UserName    string `json:"userName"`
	AvatarId    string `json:"avatarId"`
	Title       string `json:"title"`
	Discription string `json:"discription"`
	Intrest     string `json:"interest"`
	ThumbnailId string `json:"thumbnailId"`
}
