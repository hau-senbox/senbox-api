package entity

import (
	"sen-global-api/internal/domain/value"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SStaffFormApplication struct {
	ID             uuid.UUID                   `gorm:"column:id;type:char(36);primary_key"`
	UserID         uuid.UUID                   `gorm:"column:user_id;type:char(36);not null"`
	OrganizationID uuid.UUID                   `gorm:"column:organization_id;type:char(36);not null"`
	Status         value.FromApplicationStatus `gorm:"column:status;not null"`
	ApprovedAt     time.Time                   `gorm:"column:approved_at;type:datetime"`
	CreatedAt      time.Time                   `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func (application *SStaffFormApplication) BeforeCreate(tx *gorm.DB) (err error) {
	application.Status = value.Pending

	return err
}

func (t *SStaffFormApplication) IsInOrganizations(orgIDs []string) bool {
	orgIDStr := t.OrganizationID.String()
	for _, id := range orgIDs {
		if orgIDStr == id {
			return true
		}
	}
	return false
}
