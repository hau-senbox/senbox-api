package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sen-global-api/internal/domain/value"
	"time"
)

type SStudentFormApplication struct {
	ID             int64                       `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	StudentName    string                      `gorm:"column:student_name;type:varchar(255);not null"`
	UserID         uuid.UUID                   `gorm:"column:user_id;primary_key"`
	User           SUserEntity                 `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE;"`
	OrganizationID uuid.UUID                   `gorm:"column:organization_id;primary_key"`
	Organization   SOrganization               `gorm:"foreignKey:OrganizationID;references:id;constraint:OnDelete:CASCADE;"`
	Status         value.FromApplicationStatus `gorm:"column:status;not null"`
	ApprovedAt     time.Time                   `gorm:"column:approved_at;type:datetime"`
	CreatedAt      time.Time                   `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func (application *SStudentFormApplication) BeforeCreate(tx *gorm.DB) (err error) {
	application.Status = value.Pending

	return err
}
