package domain

type Video struct {
	ID          uint   `gorm:"primarykey;auto_increment"`
	VideoId     string `gorm:"not null"`
	Archived    bool   `gorm:"default:false"`
	S3Path      string `json:"s3path"`
	UserName    string `json:"userName"`
	AvatarId    string `json:"avatarId"`
	Title       string `json:"title"`
	Discription string `json:"discription"`
	Interest    string `json:"interest"`
	ThumbnailId string `json:"thumbnailId"`
}

type ToSaveVideo struct {
	S3Path      string `json:"s3path"`
	UserName    string `json:"userName"`
	AvatarId    string `json:"avatarId"`
	Title       string `json:"title"`
	Discription string `json:"discription"`
	Interest    string `json:"interest"`
	ThumbnailId string `json:"thumbnailId"`
}
