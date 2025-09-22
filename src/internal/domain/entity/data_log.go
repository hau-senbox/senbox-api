package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
)

type DataLog struct {
	ID           uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Endpoint     string         `gorm:"type:varchar(255);not null" json:"endpoint"`
	Method       string         `gorm:"type:varchar(10);not null" json:"method"`
	Payload      datatypes.JSON `gorm:"type:json;not null" json:"payload"`
	Response     datatypes.JSON `gorm:"type:json" json:"response,omitempty"`
	Status       string         `gorm:"type:varchar(20);not null" json:"status"`
	ErrorMessage *string        `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
}
