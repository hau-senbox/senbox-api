package entity

import (
	"sen-global-api/internal/domain/value"
	"time"

	"gorm.io/datatypes"
)

type SyncQueue struct {
	ID               uint64                `gorm:"primaryKey;autoIncrement" json:"id"`
	LastSubmissionID uint64                `gorm:"not null" json:"last_submission_id"`
	LastSubmittedAt  time.Time             `gorm:"not null" json:"last_submitted_at"`
	FormNotes        datatypes.JSON        `gorm:"type:json" json:"form_notes"`
	SheetName        string                `gorm:"type:varchar(128);not null" json:"sheet_name"`
	SpreadsheetID    string                `gorm:"type:varchar(128);not null" json:"spreadsheet_id"`
	SheetUrl         string                `gorm:"type:varchar(255);not null" json:"sheet_url"`
	Status           value.SyncQueueStatus `gorm:"type:varchar(32);not null" json:"status"`
	CreatedAt        time.Time             `json:"created_at"`
	UpdatedAt        time.Time             `json:"updated_at"`
}
