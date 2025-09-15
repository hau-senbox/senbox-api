package entity

import (
	"time"

	"github.com/google/uuid"
)

type SuperAdminEmergencyMenu struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	ComponentID uuid.UUID `json:"component_id" gorm:"type:char(36);not null"`
	Order       int       `json:"order" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
