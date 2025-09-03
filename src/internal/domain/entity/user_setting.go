package entity

import (
	"time"

	"gorm.io/datatypes"
)

type UserSetting struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID   string         `gorm:"type:varchar(255);not null;default:''"`
	OwnerRole string         `gorm:"type:varchar(255);not null;default:''"`
	Key       string         `gorm:"unique;not null" json:"key"`
	Value     datatypes.JSON `gorm:"type:json" json:"value"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
