package entity

import (
	"time"

	"gorm.io/datatypes"
)

type SDeviceComponentValues struct {
	ID             int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Setting        datatypes.JSON `gorm:"column:setting;type:json;not null;default:'{}'"`
	OrganizationId *int64         `gorm:"column:organization_id;"`
	Organization   *SOrganization `gorm:"foreignKey:OrganizationId;references:id;constraint:OnDelete:CASCADE;default:1"`
	CreatedAt      time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}
