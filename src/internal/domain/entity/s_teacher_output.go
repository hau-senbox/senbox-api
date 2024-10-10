package entity

import "time"

type STeacherOutput struct {
	Value2        string    `gorm:"column:value2;primary_key;not null;index" json:"value2"`
	Value3        string    `gorm:"column:value3;primary_key;not null;index" json:"value3"`
	SpreadsheetID string    `gorm:"column:spreadsheet_id;not null;"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"updated_at"`
}
