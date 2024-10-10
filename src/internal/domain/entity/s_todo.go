package entity

import (
	"gorm.io/datatypes"
	"sen-global-api/internal/domain/value"
	"time"
)

type Task struct {
	Index     int    `json:"index"`
	Name      string `json:"name"`
	DueDate   string `json:"due_date"`
	Value     string `json:"value"`
	Selection string `json:"selection"`
	Selected  string `json:"selected"`
}

type STasks struct {
	Tasks []Task `json:"tasks"`
}

type SToDo struct {
	ID                   string                     `gorm:"primary_key;type:varchar(255);not null" json:"id"`
	Name                 string                     `gorm:"type:varchar(255);" json:"name"`
	Type                 value.ToDoType             `gorm:"type:varchar(32);default:'assign'" json:"type"`
	SpreadsheetID        string                     `gorm:"type:varchar(255);not null" json:"spreadsheet_id"`
	SheetName            string                     `gorm:"type:varchar(255);not null;default:Tasks" json:"sheet_name"`
	Tasks                datatypes.JSONType[STasks] `gorm:"type:json;not null;default:'[]'" json:"tasks"`
	HistorySpreadsheetID string                     `gorm:"type:varchar(255);not null" json:"history_spreadsheet_id"`
	HistorySheetName     string                     `gorm:"type:varchar(255);not null;default:Answers" json:"history_sheet_name"`
	StartRow             int                        `gorm:"type:int;not null;default:13" json:"start_row"`
	CreatedAt            time.Time                  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt            time.Time                  `gorm:"type:timestamp;not null;default:CURRENT_TIMESTAMP" json:"updated_at"`
}
