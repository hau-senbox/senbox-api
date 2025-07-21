package entity

import (
	"time"

	"github.com/google/uuid"
)

type StudentMenu struct {
	ID          uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	StudentID   uuid.UUID `json:"student_id" gorm:"type:char(36);primaryKey"`
	ComponentID uuid.UUID `json:"component_id" gorm:"type:char(36);primaryKey"`
	Order       int       `json:"order" gorm:"not null"`
	IsShow      bool      `json:"is_show" gorm:"default:true"`
	Visible     bool      `json:"visible" gorm:"default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
