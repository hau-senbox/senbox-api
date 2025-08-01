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
}

func (application *SStudentFormApplication) BeforeCreate(tx *gorm.DB) (err error) {
	application.Status = value.Pending

	return err
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
