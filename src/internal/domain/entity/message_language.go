package entity

import "time"

type MessageLanguage struct {
	ID         int       `gorm:"primaryKey;autoIncrement"`
	TypeID     string    `gorm:"type:varchar(50);not null;default:''"`
	Type       string    `gorm:"type:varchar(50);not null;default:''"`
	Key        string    `gorm:"type:varchar(100);not null;default:''"`
	Value      string    `gorm:"type:varchar(255);not null;default:''"`
	LanguageID uint      `gorm:"not null;default:1"`
	CreatedAt  time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt  time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
