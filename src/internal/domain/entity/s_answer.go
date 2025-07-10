package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type SAnswer struct {
	ID           uuid.UUID       `gorm:"type:char(36);primary_key"`
	UserID       string          `gorm:"type:varchar(255);not null;default:''"`
	SubmissionID uint64          `gorm:"not null"`
	Response     json.RawMessage `gorm:"type:json" json:"response"`
	Key          string          `gorm:"type:varchar(255);not null;default:''"`
	DB           string          `gorm:"type:varchar(255);not null;default:''"`
	CreatedAt    time.Time       `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt    time.Time       `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
