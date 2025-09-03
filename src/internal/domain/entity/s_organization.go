package entity

import (
	"time"

	"github.com/google/uuid"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SOrganization struct {
	ID                   uuid.UUID `gorm:"type:char(36);primary_key"`
	OrganizationName     string    `gorm:"type:varchar(255);not null;unique"`
	OrganizationNickName string    `gorm:"type:varchar(255);not null;default:''"`
	Avatar               string    `gorm:"type:varchar(255);not null;default:''"`
	AvatarURL            string    `gorm:"type:longtext;not null;default:''"`
	Password             string    `gorm:"type:varchar(255);not null;default:''"`
	Address              string    `gorm:"type:varchar(255);not null;default:''"`
	Description          string    `gorm:"type:varchar(255);not null;default:''"`
	CreatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt            time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	CreatedIndex         int       `gorm:"column:created_index;not null;default:0"`

	UserOrgs []SUserOrg `gorm:"foreignKey:organization_id;references:id;constraint:OnDelete:CASCADE"`
}

func (organization *SOrganization) BeforeCreate(tx *gorm.DB) (err error) {
	// Tạo UUID mới
	id, err := uuid.NewUUID()
	if err == nil {
		organization.ID = id
	}

	// Hash password
	encryptedPwdData, err := bcrypt.GenerateFromPassword([]byte(organization.Password), bcrypt.DefaultCost)
	if err == nil {
		organization.Password = string(encryptedPwdData)
	}

	// Tính CreatedIndex = MAX(created_index) + 1
	var maxIndex int
	if err := tx.Model(&SOrganization{}).
		Select("COALESCE(MAX(created_index), 0)").
		Scan(&maxIndex).Error; err != nil {
		return err
	}
	organization.CreatedIndex = maxIndex + 1

	return nil
}
