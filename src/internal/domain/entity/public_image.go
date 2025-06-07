package entity

type PublicImage struct {
	ID        uint64 `gorm:"primary_key;auto_increment;"`
	ImageName string `gorm:"column:image_name;not null;"`
	Folder    string `gorm:"column:folder;not null;"`
	Key       string `gorm:"column:key;not null;unique;"`
	URL       string `gorm:"column:url;not null;"`
	Extension string `gorm:"column:extension;not null;"`
	Width     int    `gorm:"column:width;not null;default:0;"`
	Height    int    `gorm:"column:height;not null;default:0;"`
}
