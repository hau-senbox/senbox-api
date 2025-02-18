package entity

import (
	"database/sql"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type SUserEntity struct {
	ID           uuid.UUID     `gorm:"type:char(36);primary_key"`
	Username     string        `gorm:"type:varchar(255);not null;default:''"`
	Fullname     string        `gorm:"type:varchar(255);not null;default:''"`
	Phone        string        `gorm:"type:varchar(255);not null;default:''"`
	Email        string        `gorm:"type:varchar(255);not null;default:''"`
	Birthday     time.Time     `gorm:"default:CURRENT_TIMESTAMP"`
	Password     string        `gorm:"type:varchar(255);not null;default:''"`
	UserConfigID sql.NullInt64 `gorm:"column:user_config_id"`
	UserConfig   *SUserConfig  `gorm:"foreignKey:UserConfigID;references:id;constraint:OnDelete:CASCADE"`
	CompanyId    int64         `gorm:"column:company_id;"`
	Company      SCompany      `gorm:"foreignKey:CompanyId;references:id;constraint:OnDelete:CASCADE;default:1"`
	CreatedAt    time.Time     `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt    time.Time     `gorm:"default:CURRENT_TIMESTAMP;not null"`

	// Many-to-many relationship with roles
	Roles        []SRole       `gorm:"many2many:s_user_roles;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:role_id"`
	RolePolicies []SRolePolicy `gorm:"many2many:s_user_policies;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:policy_id"`

	Guardians []SUserEntity `gorm:"many2many:s_user_guardians;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:guardian_id"`
	Devices   []SDevice     `gorm:"many2many:s_user_devices;foreignKey:id;joinForeignKey:user_id;references:id;joinReferences:device_id"`
}

func (user *SUserEntity) BeforeCreate(tx *gorm.DB) (err error) {
	id, err := uuid.NewUUID()
	if err == nil {
		user.ID = id
	}

	if user.Password != "" {
		encryptedPwdData, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err == nil {
			user.Password = string(encryptedPwdData)
		}
	}

	user.Username = strings.ToLower(html.EscapeString(strings.TrimSpace(user.Username)))

	return err
}
