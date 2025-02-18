package entity

import (
	"time"

	"gorm.io/datatypes"
)

type SDeviceComponentValues struct {
	ID        int64          `gorm:"column:id;primaryKey;autoIncrement"`
	Setting   datatypes.JSON `gorm:"column:setting;type:json;not null;default:'{}'"`
	CompanyId *int64         `gorm:"column:company_id;"`
	Company   *SCompany      `gorm:"foreignKey:CompanyId;references:id;constraint:OnDelete:CASCADE;default:1"`
	CreatedAt time.Time      `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt time.Time      `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}
