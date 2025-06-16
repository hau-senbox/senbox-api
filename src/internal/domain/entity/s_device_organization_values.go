package entity

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/datatypes"
)

type SDeviceComponentValues struct {
	ID             int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Setting        datatypes.JSON `gorm:"column:setting;type:json;not null;default:'{}'"`
	OrganizationID *uuid.UUID     `gorm:"column:organization_id;"`
	Organization   *SOrganization `gorm:"foreignKey:OrganizationID;references:id;constraint:OnDelete:CASCADE;default:1"`
	CreatedAt      time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}
