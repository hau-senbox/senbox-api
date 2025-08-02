package entity

import (
	"sen-global-api/internal/domain/value"
	"time"

	"gorm.io/datatypes"
)

type SyncQueue struct {
	ID               uint64                `gorm:"primaryKey;autoIncrement"`
	LastSubmissionID uint64                `gorm:"not null"`
	LastSubmittedAt  string                `gorm:"not null"`
	FormNotes        datatypes.JSON        `gorm:"type:json"`
	SheetName        string                `gorm:"type:varchar(128);not null"`
	SpreadsheetID    string                `gorm:"type:varchar(128);not null"`
	Status           value.SyncQueueStatus `gorm:"type:varchar(32);not null"` // pending, done, failed
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
