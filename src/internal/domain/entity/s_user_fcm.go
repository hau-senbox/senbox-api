package entity

import "time"

type SUserFCMToken struct {
	ID        uint64    `gorm:"primary_key;auto_increment;not null"`
	UserID    string    `gorm:"type:varchar(255);not null"`
	FCMToken  string    `gorm:"type:varchar(255);not null"`
	DeviceID  string    `gorm:"type:varchar(255);not null"`
	IsActive  bool      `gorm:"type:tinyint;not null;default:1"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func (*SUserFCMToken) TableName() string {
	return "s_user_fcm"
}