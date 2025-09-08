package entity

import (
	"time"

	"github.com/google/uuid"
)

type DepartmentMenuOrganization struct {
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	DepartmentID   string    `json:"department_id" gorm:"type:varchar(255);not null"`
	OrganizationID string    `json:"organization_id" gorm:"type:varchar(255);not null"`
	ComponentID    uuid.UUID `json:"component_id" gorm:"type:char(36);not null"`
	Order          int       `json:"order" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
