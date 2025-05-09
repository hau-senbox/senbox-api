package entity

type SImage struct {
	ID        uint64 `gorm:"primary_key;auto_increment;"`
	ImageName string `gorm:"column:image_name;not null;"`
	Folder    string `gorm:"column:folder;not null;"`
	Key       string `gorm:"column:key;not null;unique;"`
	Extension string `gorm:"column:extension;not null;"`
	Width     int    `gorm:"column:width;not null;default:0;"`
	Height    int    `gorm:"column:height;not null;default:0;"`
}

func (SImage) TableName() string {
	return "s_image"
}
