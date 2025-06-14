package entity

type SOrgDevices struct {
	OrganizationID int64         `gorm:"column:organization_id;primary_key"`
	Organization   SOrganization `gorm:"foreignKey:OrganizationID;references:id;constraint:OnDelete:CASCADE;"`
	DeviceID       string        `gorm:"column:device_id;primary_key"`
	Device         SDevice       `gorm:"foreignKey:DeviceID;references:id;constraint:OnDelete:CASCADE;"`
}
