package entity

type SVideo struct {
	ID        uint64 `gorm:"primary_key;auto_increment;"`
	VideoName string `gorm:"column:video_name;not null;"`
	Folder    string `gorm:"column:folder;not null;"`
	Key       string `gorm:"column:key;not null;unique;"`
	Extension string `gorm:"column:extension;not null;"`
}
