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
	IsAdminBlock   bool                        `gorm:"column:is_admin_block;default:false"`
	ApprovedAt     time.Time                   `gorm:"column:approved_at;type:datetime"`
	CreatedAt      time.Time                   `gorm:"default:CURRENT_TIMESTAMP;not null"`
	CreatedIndex   int                         `gorm:"column:created_index;not null;default:0"`
}

func (application *SStaffFormApplication) BeforeCreate(tx *gorm.DB) (err error) {
	application.Status = value.Pending

	// Tính CreatedIndex = MAX(created_index) + 1 theo OrganizationID
	var maxIndex int
	if err := tx.Model(&SStaffFormApplication{}).
		Where("organization_id = ?", application.OrganizationID).
		Select("COALESCE(MAX(created_index), 0)").
		Scan(&maxIndex).Error; err != nil {
		return err
	}
	application.CreatedIndex = maxIndex + 1

	return nil
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
