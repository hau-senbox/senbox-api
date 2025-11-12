package entity

import (
	"sen-global-api/internal/domain/value"
	"time"

	"gorm.io/gorm"
)

type SDevice struct {
	ID                      string                 `gorm:"type:varchar(36);primary_key;"`
	DeviceName              string                 `gorm:"type:varchar(255);not null;default:''"`
	InputMode               value.InfoInputType    `gorm:"type:varchar(32);not null;default:1"`
	ScreenButtonType        value.ScreenButtonType `gorm:"type:varchar(16);not null;default:'scan'"`
	Status                  value.DeviceMode       `gorm:"type:varchar(32);not null;default:1"`
	DeactivateMessage       string                 `gorm:"type:varchar(255);not null;default:''"`
	ButtonUrl               string                 `gorm:"type:varchar(255);not null;default:''"`
	Note                    string                 `gorm:"type:varchar(255);not null;default:''"`
	AppVersion              string                 `gorm:"type:varchar(255);not null;default:''"`
	RowNo                   int                    `gorm:"type:int;not null;default:0"`
	DeviceComponentValuesID int64                  `gorm:"column:device_component_values_id;default:1"`
	DeviceComponentValues   SDeviceComponentValues `gorm:"foreignKey:DeviceComponentValuesID;references:id;constraint:OnDelete:CASCADE"`
	CreatedAt               time.Time              `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt               time.Time              `gorm:"default:CURRENT_TIMESTAMP;not null"`
	CreatedIndex            int                    `gorm:"type:int;not null;default:0"`
}

func (d *SDevice) BeforeCreate(tx *gorm.DB) error {
	var maxIndex int
	if err := tx.Model(&SDevice{}).
		Select("COALESCE(MAX(created_index), 0)").
		Scan(&maxIndex).Error; err != nil {
		return err
	}
	d.CreatedIndex = maxIndex + 1
	return nil
}
