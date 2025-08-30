package entity

import "time"

type UserDevicesLogin struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	UserID   string    `json:"user_id" gorm:"type:varchar(255);not null;index:user_device_login_idx,priority:1"`
	DeviceID string    `json:"device_id" gorm:"type:varchar(255);not null;index:user_device_login_idx,priority:2"`
	LoginAt  time.Time `json:"login_at" gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP"`
}
