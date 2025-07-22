package entity

type SAudio struct {
	ID        uint64 `gorm:"primary_key;auto_increment;"`
	AudioName string `gorm:"column:audio_name;not null;"`
	Folder    string `gorm:"column:folder;not null;"`
	Key       string `gorm:"column:key;not null;unique;"`
	Extension string `gorm:"column:extension;not null;"`
}
