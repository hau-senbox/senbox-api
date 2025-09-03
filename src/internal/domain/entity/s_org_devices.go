package entity

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SOrgDevices struct {
	OrganizationID uuid.UUID     `gorm:"column:organization_id;primary_key"`
	Organization   SOrganization `gorm:"foreignKey:OrganizationID;references:id;constraint:OnDelete:CASCADE;"`
	DeviceID       string        `gorm:"column:device_id;primary_key"`
	Device         SDevice       `gorm:"foreignKey:DeviceID;references:id;constraint:OnDelete:CASCADE;"`
	DeviceName     string        `gorm:"column:device_name;type:varchar(255);default:''"`
	DeviceNickName string        `gorm:"column:device_nick_name;type:varchar(255);default:''"`
	CreatedIndex   int           `gorm:"column:created_index;not null;default:0"`
}

func (d *SOrgDevices) BeforeCreate(tx *gorm.DB) (err error) {
	var maxIndex int
	if err = tx.Model(&SOrgDevices{}).
		Where("organization_id = ?", d.OrganizationID).
		Select("COALESCE(MAX(created_index), 0)").
		Scan(&maxIndex).Error; err != nil {
		return err
	}

	// Gán created_index = max + 1
	d.CreatedIndex = maxIndex + 1
	return nil
}
