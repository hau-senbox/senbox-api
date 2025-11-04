package entity

import (
	"sen-global-api/internal/domain/value"
	"time"

	"gorm.io/datatypes"
)

type UserSetting struct {
	ID        uint                 `gorm:"primaryKey;autoIncrement" json:"id"`
	OwnerID   string               `gorm:"type:varchar(255);not null;default:''"`
	OwnerRole value.OwnerRole      `gorm:"type:varchar(255);not null;default:''"`
	Key       value.UserSettingKey `gorm:"not null" json:"key"`
	Value     datatypes.JSON       `gorm:"type:json" json:"value"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserSettingWelcomeReminder struct {
	IsEnabled    bool   `json:"is_enabled" binding:"required"`
	TimeReminder string `json:"time_reminder"`
}
