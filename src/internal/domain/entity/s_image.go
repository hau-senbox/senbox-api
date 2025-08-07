package entity

import "time"

type SImage struct {
	ID        uint64    `gorm:"primary_key;auto_increment;"`
	ImageName string    `gorm:"column:image_name;not null;"`
	Folder    string    `gorm:"column:folder;not null;"`
	Key       string    `gorm:"column:key;not null;unique;"`
	Extension string    `gorm:"column:extension;not null;"`
	Width     int       `gorm:"column:width;not null;default:0;"`
	Height    int       `gorm:"column:height;not null;default:0;"`
	TopicID   string    `gorm:"column:topic_id;not null;default:''"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime"`
}

func (SImage) TableName() string {
	return "s_image"
}
