package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StudentMenuOrganization struct {
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	StudentID      string    `json:"student_id" gorm:"type:char(36);not null"`
	OrganizationID string    `json:"organization_id" gorm:"type:char(36);not null"`
	ComponentID    string    `json:"component_id" gorm:"type:char(36);not null"`
	Order          int       `json:"order" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (s *StudentMenuOrganization) BeforeCreate(tx *gorm.DB) (err error) {
	if s.ID == uuid.Nil {
		s.ID = uuid.New()
	}
	return
}
