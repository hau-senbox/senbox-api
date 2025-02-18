package entity

import (
	"github.com/google/uuid"
)

type SUserDevices struct {
	UserId   uuid.UUID   `gorm:"column:user_id;primary_key"`
	User     SUserEntity `gorm:"foreignKey:UserId;references:id;constraint:OnDelete:CASCADE;"`
	DeviceId string      `gorm:"column:device_id;primary_key"`
	Device   SDevice     `gorm:"foreignKey:DeviceId;references:id;constraint:OnDelete:CASCADE;"`
}
