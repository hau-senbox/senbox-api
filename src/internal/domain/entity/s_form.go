package entity

import (
	"sen-global-api/internal/domain/value"
	"strings"
	"time"

	"gorm.io/gorm"
)

type SForm struct {
	ID             uint64         `gorm:"primary_key;AUTO_INCREMENT"`
	Note           string         `gorm:"type:varchar(255);not null;unique"`
	Name           string         `gorm:"type:varchar(1000);not null;default:''"`
	SpreadsheetUrl string         `gorm:"type:varchar(255);not null"`
	SpreadsheetId  string         `gorm:"type:varchar(255);not null"`
	Password       string         `gorm:"type:varchar(255);"`
	Status         value.Status   `gorm:"type:tinyint;not null;default:1"`
	SheetName      string         `gorm:"type:varchar(255);not null;default:'Questions'"`
	Type           value.FormType `gorm:"type:varchar(32);not null;default:'general'"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func (receiver *SForm) BeforeCreate(tx *gorm.DB) (err error) {
	receiver.Type = value.FormType_General
	if strings.Contains(strings.ToLower(receiver.Note), "memory-form") {
		receiver.Type = value.FormType_SelfRemember
	}

	return nil
}

func (receiver *SForm) BeforeUpdate(tx *gorm.DB) (err error) {
	receiver.Type = value.FormType_General
	if strings.Contains(strings.ToLower(receiver.Note), "memory-form") {
		receiver.Type = value.FormType_SelfRemember
	}

	return nil
}
