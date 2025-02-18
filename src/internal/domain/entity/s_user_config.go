package entity

import (
	"time"
)

type SUserConfig struct {
	ID                   int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	TopButtonConfig      string    `gorm:"type:varchar(255);not null;default:''"`
	StudentOutputSheetId string    `gorm:"type:varchar(255);not null;"`
	TeacherOutputSheetId string    `gorm:"type:varchar(255);not null;default:''"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
