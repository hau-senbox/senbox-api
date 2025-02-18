package entity

import (
	"time"
)

type SCompany struct {
	ID          int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	CompanyName string    `gorm:"type:varchar(255);not null;"`
	Address     string    `gorm:"type:varchar(255);not null;default:'';unique;unique_index"`
	Description string    `gorm:"type:varchar(255);not null;default:''"`
	CreatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt   time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
