package entity

import (
	"sen-global-api/internal/domain/value"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SStudentFormApplication struct {
	ID             uuid.UUID                   `gorm:"column:id;type:char(36);primary_key"`
	StudentName    string                      `gorm:"column:student_name;type:varchar(255);not null"`
	ChildID        uuid.UUID                   `gorm:"column:child_id;type:char(36);not null"`
	UserID         uuid.UUID                   `gorm:"column:user_id;type:char(36);not null"`
	CustomID       string                      `gorm:"column:custom_id;type:varchar(255);not null;default:''"`
	OrganizationID uuid.UUID                   `gorm:"column:organization_id;type:char(36);not null"`
	Status         value.FromApplicationStatus `gorm:"column:status;not null"`
	IsAdminBlock   bool                        `gorm:"column:is_admin_block;default:false"`
	ApprovedAt     time.Time                   `gorm:"column:approved_at;type:datetime"`
	CreatedAt      time.Time                   `gorm:"default:CURRENT_TIMESTAMP;not null"`
	CreatedIndex   int                         `gorm:"column:created_index;not null;default:0"`
}

func (application *SStudentFormApplication) BeforeCreate(tx *gorm.DB) (err error) {
	// Luôn set trạng thái Pending khi tạo mới
	application.Status = value.Pending

	// Tính CreatedIndex = MAX(created_index) + 1 theo organization_id
	var maxIndex int
	if err := tx.Model(&SStudentFormApplication{}).
		Where("organization_id = ?", application.OrganizationID).
		Select("COALESCE(MAX(created_index), 0)").
		Scan(&maxIndex).Error; err != nil {
		return err
	}

	application.CreatedIndex = maxIndex + 1

	return nil
}

func (application *SStudentFormApplication) IsInOrganizations(orgIDs []string) bool {
	orgIDStr := application.OrganizationID.String()

	for _, id := range orgIDs {
		if id == orgIDStr {
			return true
		}
	}

	return false
}
