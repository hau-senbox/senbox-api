package entity

import (
	"time"
)

type SRole struct {
	ID        int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	RoleName  string    `gorm:"type:varchar(255);not null;default:''"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
