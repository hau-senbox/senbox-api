package entity

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SOrganization struct {
	ID               int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	OrganizationName string    `gorm:"type:varchar(255);not null;"`
	Password         string    `gorm:"type:varchar(255);not null;default:''"`
	Address          string    `gorm:"type:varchar(255);not null;default:''"`
	Description      string    `gorm:"type:varchar(255);not null;default:''"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func (organization *SOrganization) BeforeCreate(tx *gorm.DB) (err error) {
	encryptedPwdData, err := bcrypt.GenerateFromPassword([]byte(organization.Password), bcrypt.DefaultCost)
	if err == nil {
		organization.Password = string(encryptedPwdData)
	}

	return err
}
