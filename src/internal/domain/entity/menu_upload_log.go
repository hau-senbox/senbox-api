package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type MenuUploadLog struct {
	ID           uuid.UUID      `gorm:"type:char(36);primaryKey" json:"id"`
	Endpoint     string         `gorm:"type:varchar(255);not null" json:"endpoint"`
	Method       string         `gorm:"type:varchar(10);not null" json:"method"`
	Payload      datatypes.JSON `gorm:"type:json;not null" json:"payload"`
	Status       string         `gorm:"type:varchar(20);not null" json:"status"`
	ErrorMessage *string        `gorm:"type:text" json:"error_message,omitempty"`
	CreatedAt    time.Time      `gorm:"type:datetime(3);not null;default:CURRENT_TIMESTAMP(3)" json:"created_at"`
}

// Hook tạo UUID tự động trước khi insert
func (m *MenuUploadLog) BeforeCreate(tx *gorm.DB) (err error) {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return
}
