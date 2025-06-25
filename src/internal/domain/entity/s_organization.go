package entity

import (
	"github.com/google/uuid"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SOrganization struct {
	ID               uuid.UUID `gorm:"type:char(36);primary_key"`
	OrganizationName string    `gorm:"type:varchar(255);not null;unique"`
	Avatar           string    `gorm:"type:varchar(255);not null;default:''"`
	AvatarURL        string    `gorm:"type:longtext;not null;default:''"`
	Password         string    `gorm:"type:varchar(255);not null;default:''"`
	Address          string    `gorm:"type:varchar(255);not null;default:''"`
	Description      string    `gorm:"type:varchar(255);not null;default:''"`
	CreatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt        time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`

	UserOrgs []SUserOrg `gorm:"foreignKey:organization_id;references:id;constraint:OnDelete:CASCADE"`
}

func (organization *SOrganization) BeforeCreate(tx *gorm.DB) (err error) {
	encryptedPwdData, err := bcrypt.GenerateFromPassword([]byte(organization.Password), bcrypt.DefaultCost)
	if err == nil {
		organization.Password = string(encryptedPwdData)
	}

	return err
}
