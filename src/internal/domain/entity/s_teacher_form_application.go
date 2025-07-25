package entity

import (
	"sen-global-api/internal/domain/value"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type STeacherFormApplication struct {
	ID             uuid.UUID                   `gorm:"column:id;type:char(36);primary_key"`
	UserID         uuid.UUID                   `gorm:"column:user_id;type:char(36);not null"`
	OrganizationID uuid.UUID                   `gorm:"column:organization_id;type:char(36);not null"`
	Status         value.FromApplicationStatus `gorm:"column:status;not null"`
	ApprovedAt     time.Time                   `gorm:"column:approved_at;type:datetime"`
	CreatedAt      time.Time                   `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func (application *STeacherFormApplication) BeforeCreate(tx *gorm.DB) (err error) {
	application.Status = value.Pending

	return err
}
