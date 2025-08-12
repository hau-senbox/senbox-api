package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationSetting struct {
	ID                uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`
	OrganizationID    string    `gorm:"type:varchar(255);not null" json:"organization_id"`
	ComponentID       string    `gorm:"type:varchar(255);not null" json:"component_id"`
	IsViewMessage     bool      `gorm:"not null;default:false" json:"is_view_message"`
	IsShowOrgNews     bool      `gorm:"not null;default:false" json:"is_show_org_news"`
	IsDeactiveTopMenu bool      `gorm:"not null;default:false" json:"is_deactive_top_menu"`
	IsShowSpecialBtn  bool      `gorm:"not null;default:false" json:"is_show_special_btn"`
	MessageBox        string    `gorm:"type:text" json:"message_box"`
	MessageTopMenu    string    `gorm:"type:text" json:"message_top_menu"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (o *OrganizationSetting) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
