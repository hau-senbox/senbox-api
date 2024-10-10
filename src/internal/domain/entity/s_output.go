package entity

import "time"

type SOutput struct {
	Value1        string    `gorm:"column:value1;primary_key;not null;index" json:"value1"`
	Value2        string    `gorm:"column:value2;primary_key;not null;index" json:"value2"`
	SpreadsheetID string    `gorm:"column:spreadsheet_id;not null;unique;uniqueIndex" json:"spreadsheet_id"`
	CreatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"created_at"`
	UpdatedAt     time.Time `gorm:"default:CURRENT_TIMESTAMP;not null" json:"updated_at"`
}
