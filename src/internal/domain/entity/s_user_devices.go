package entity

import (
	"github.com/google/uuid"
)

type SUserDevices struct {
	UserID   uuid.UUID   `gorm:"column:user_id;primary_key"`
	User     SUserEntity `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE;"`
	DeviceID string      `gorm:"column:device_id;primary_key"`
	Device   SDevice     `gorm:"foreignKey:DeviceID;references:id;constraint:OnDelete:CASCADE;"`
}
