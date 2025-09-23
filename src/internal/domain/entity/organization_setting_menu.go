package entity

import (
	"time"

	"github.com/google/uuid"
)

type OrganizationSettingMenu struct {
	ID                    uuid.UUID `json:"id" gorm:"type:char(36);primaryKey"`
	OrganizationSettingID string    `json:"organization_setting_id" gorm:"type:varchar(255);not null"`
	ComponentID           uuid.UUID `json:"component_id" gorm:"type:char(36);not null"`
	CreatedAt             time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt             time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}
