package entity

import (
	"database/sql"
	"time"
)

type SCodeCounting struct {
	ID           uint         `gorm:"primarykey;autoIncrement" json:"id"`
	CreatedAt    time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt    time.Time    `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	DeletedAt    sql.NullTime `gorm:"index" json:"deleted_at"`
	Token        string       `gorm:"type:varchar(255);primary_key;not null" json:"token" binding:"required"`
	CurrentValue int          `gorm:"type:int;not null;default:0" json:"current_value" binding:"required"`
}
