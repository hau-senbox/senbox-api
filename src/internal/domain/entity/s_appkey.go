package entity

import "time"

type SAppKey struct {
	ID        int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	AppKey    string    `gorm:"column:app_key;type:varchar(255);not null"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
