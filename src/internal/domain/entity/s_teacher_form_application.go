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
	IsAdminBlock   bool                        `gorm:"column:is_admin_block;default:false"`
	ApprovedAt     time.Time                   `gorm:"column:approved_at;type:datetime"`
	CreatedAt      time.Time                   `gorm:"default:CURRENT_TIMESTAMP;not null"`
	CreatedIndex   int                         `gorm:"column:created_index;not null;default:0"`
}

func (application *STeacherFormApplication) BeforeCreate(tx *gorm.DB) (err error) {
	// Default status
	application.Status = value.Pending

	// Đếm số record hiện có theo OrganizationID
	var count int64
	err = tx.Model(&STeacherFormApplication{}).
		Where("organization_id = ?", application.OrganizationID).
		Count(&count).Error
	if err != nil {
		return err
	}

	application.CreatedIndex = int(count) + 1

	return nil
}

func (t *STeacherFormApplication) IsInOrganizations(orgIDs []string) bool {
	orgIDStr := t.OrganizationID.String()
	for _, id := range orgIDs {
		if orgIDStr == id {
			return true
		}
	}
	return false
}
