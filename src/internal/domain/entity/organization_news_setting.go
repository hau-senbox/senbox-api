package entity

import "github.com/google/uuid"

type OrganizationNewsSetting struct {
	ID                 uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	OrganizationID     string    `json:"organization_id"`
	IsPusblishedDevice bool      `json:"is_pusblished_device"`
	MessageDeviceNews  string    `json:"message_devices_news" binding:"required"`
	IsPusblishedPortal bool      `json:"is_pusblished_portal"`
	MessagePortalNews  string    `json:"message_portal_news" binding:"required"`
}
