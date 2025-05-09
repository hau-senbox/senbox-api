package entity

import "time"

type SRedirectUrl struct {
	ID           uint64    `gorm:"primary_key;auto_increment;not null"`
	QRCode       string    `gorm:"type:varchar(255);not null;unique"`
	TargetUrl    string    `gorm:"type:varchar(255);not null"`
	Password     *string   `gorm:"type:varchar(32);"`
	Hint         string    `gorm:"type:varchar(255);default:''"`
	HashPassword *string   `gorm:"type:varchar(255);default:null"`
	CreatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt    time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
