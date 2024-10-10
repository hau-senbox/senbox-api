package entity

import (
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
)

type SMobileDevice struct {
	gorm.Model
	DeviceId string           `gorm:"type:varchar(36);not null;unique"`
	Type     value.DeviceType `gorm:"type:varchar(16);not null"`
	FCMToken string           `gorm:"type:mediumtext;not null"`
}
