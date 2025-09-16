package entity

import "time"

type LanguageSetting struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	LangKey   string    `gorm:"unique;not null" json:"lang_key"`
	RegionKey string    `gorm:"type:varchar(255);not null" json:"region_key"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
}
