package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeacherMenuOrganization struct {
	ID             uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	TeacherID      string    `json:"teacher_id" gorm:"type:char(36);not null"`
	OrganizationID string    `json:"organization_id" gorm:"type:char(36);not null"`
	ComponentID    string    `json:"component_id" gorm:"type:char(36);not null"`
	Order          int       `json:"order" gorm:"not null"`
	CreatedAt      time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (t *TeacherMenuOrganization) BeforeCreate(tx *gorm.DB) (err error) {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return
}
