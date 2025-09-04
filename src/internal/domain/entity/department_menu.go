package entity

import (
	"time"

	"github.com/google/uuid"
)

type DepartmentMenu struct {
	ID           uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	DepartmentID uuid.UUID `json:"department_id" gorm:"type:char(36);not null"`
	ComponentID  uuid.UUID `json:"component_id" gorm:"type:char(36);not null"`
	Order        int       `json:"order" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
