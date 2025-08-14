package entity

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type OrganizationSetting struct {
	ID uuid.UUID `gorm:"type:char(36);primary_key" json:"id"`

	// quan ly theo device va org
	OrganizationID string `gorm:"type:varchar(255);not null" json:"organization_id"`
	DeviceID       string `gorm:"type:varchar(255);not null" json:"device_id"`

	// menu app config
	ComponentID string `gorm:"type:varchar(255);not null" json:"component_id"`

	// MessageB box app config
	IsViewMessageBox bool   `gorm:"not null;default:false" json:"is_view_message_box"`
	IsShowMessage    bool   `gorm:"not null;default:false" json:"is_show_meesage"`
	MessageBox       string `gorm:"type:text" json:"message_box"`

	// speciall btn config
	IsShowSpecialBtn bool `gorm:"not null;default:false" json:"is_show_special_btn"`

	// deactive app config
	IsDeactiveApp      bool   `gorm:"not null;default:false" json:"is_deactive_app"`
	MessageDeactiveApp string `gorm:"type:text" json:"message_deactive_app"`

	//top menu config
	IsDeactiveTopMenu bool      `gorm:"not null;default:false" json:"is_deactive_top_menu"`
	MessageTopMenu    string    `gorm:"type:text" json:"message_top_menu"`
	TopMenuPassword   string    `gorm:"type:text" json:"top_menu_password"`
	CreatedAt         time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt         time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (o *OrganizationSetting) BeforeCreate(tx *gorm.DB) (err error) {
	if o.ID == uuid.Nil {
		o.ID = uuid.New()
	}
	return
}
