package entity

import (
	"sen-global-api/internal/domain/value"
	"time"
)

type SDevice struct {
	DeviceId             string                 `gorm:"type:varchar(36);primary_key;not null"`
	DeviceName           string                 `gorm:"type:varchar(255);not null;default:''"`
	PrimaryUserInfo      string                 `gorm:"column:primary_user_info;type:varchar(255);not null;"`
	SecondaryUserInfo    string                 `gorm:"column:secondary_user_info;type:varchar(255);not null"`
	TertiaryUserInfo     string                 `gorm:"column:tertiary_user_info;type:varchar(255);not null"`
	InputMode            value.InfoInputType    `gorm:"type:varchar(32);not null;default:1"`
	ScreenButtonType     value.ScreenButtonType `gorm:"type:varchar(16);not null;default:'scan'"`
	ScreenButtonValue    string                 `gorm:"type:varchar(255);not null;default:''"`
	Status               value.DeviceMode       `gorm:"type:varchar(32);not null;default:1"`
	ProfilePictureUrl    string                 `gorm:"type:varchar(255);"`
	SpreadsheetId        string                 `gorm:"type:varchar(255);not null;"`
	TeacherSpreadsheetId string                 `gorm:"type:varchar(255);not null;default:''"`
	Message              string                 `gorm:"type:varchar(255);not null;default:''"`
	ButtonUrl            string                 `gorm:"type:varchar(255);not null;default:''"`
	Note                 string                 `gorm:"type:varchar(255);not null;default:''"`
	AppVersion           string                 `gorm:"type:varchar(255);not null;default:''"`
	RowNo                int                    `gorm:"type:int;not null;default:0"`
	CreatedAt            time.Time              `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt            time.Time              `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
