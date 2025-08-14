package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationNewsSetting struct {
	ID                 uuid.UUID `gorm:"type:char(36);primaryKey" json:"id"`
	OrganizationID     string    `gorm:"type:varchar(255);not null" json:"organization_id"`
	IsPusblishedDevice bool      `gorm:"type:tinyint(1);not null;default:0" json:"is_pusblished_device"`
	MessageDeviceNews  string    `gorm:"type:text;not null" json:"message_devices_news" binding:"required"`
	IsPusblishedPortal bool      `gorm:"type:tinyint(1);not null;default:0" json:"is_pusblished_portal"`
	MessagePortalNews  string    `gorm:"type:text;not null" json:"message_portal_news" binding:"required"`
	CreatedAt          time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdatedAt          time.Time `gorm:"type:datetime;not null;default:CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP" json:"updated_at"`
}
