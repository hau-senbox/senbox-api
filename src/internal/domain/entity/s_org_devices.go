package entity

import "github.com/google/uuid"

type SOrgDevices struct {
	OrganizationID uuid.UUID     `gorm:"column:organization_id;primary_key"`
	Organization   SOrganization `gorm:"foreignKey:OrganizationID;references:id;constraint:OnDelete:CASCADE;"`
	DeviceID       string        `gorm:"column:device_id;primary_key"`
	Device         SDevice       `gorm:"foreignKey:DeviceID;references:id;constraint:OnDelete:CASCADE;"`
	DeviceName     string        `gorm:"column:device_name;type:varchar(255);default:''"`
	DeviceNickName string        `gorm:"column:device_nick_name;type:varchar(255);default:''"`
}
