package entity

import (
	"time"

	"gorm.io/datatypes"
)

type AppConfig struct {
	ID        uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Key       string         `gorm:"unique;not null" json:"key"`
	Value     datatypes.JSON `gorm:"type:json" json:"value"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}
