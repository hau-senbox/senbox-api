package entity

import "time"

type UserBlockSetting struct {
	ID              int       `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID          string    `gorm:"type:varchar(100);not null" json:"user_id"`
	IsDeactive      bool      `gorm:"default:false" json:"is_deactive"`
	IsViewMessage   bool      `gorm:"default:true" json:"is_view_message"`
	MessageBox      string    `gorm:"type:text" json:"message_box"`
	MessageDeactive string    `gorm:"type:text" json:"message_deactive"`
	IsNeedToUpdate  bool      `gorm:"default:false" json:"is_need_to_update"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}
