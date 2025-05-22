package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"sen-global-api/internal/domain/value"
	"time"
)

type SOrgFormApplication struct {
	ID                 int64                          `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	OrganizationName   string                         `gorm:"column:organization_name;type:varchar(255);not null"`
	ApplicationContent string                         `gorm:"column:application_content;type:text;not null;default:''"`
	Status             value.OrgFromApplicationStatus `gorm:"column:status;not null"`
	ApprovedAt         time.Time                      `gorm:"column:approved_at;type:datetime"`
	CreatedAt          time.Time                      `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UserId             uuid.UUID                      `gorm:"column:user_id;primary_key"`
	User               SUserEntity                    `gorm:"foreignKey:UserId;references:id;constraint:OnDelete:CASCADE;"`
}

func (application *SOrgFormApplication) BeforeCreate(tx *gorm.DB) (err error) {
	application.Status = value.Pending

	return err
}
