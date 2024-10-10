package entity

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"html"
	"sen-global-api/internal/domain/value"
	"strings"
	"time"
)

// SUser Role: is a user role in option set
// Values: 1 Guest, 2: Device, 4: Moderator, 8: Admin, 16: Super Admin
type SUser struct {
	UserId      string     `gorm:"type:varchar(255);primary_key;not null"`
	Username    string     `gorm:"type:varchar(255);not null;default:'';unique;unique_index"`
	Fullname    string     `gorm:"type:varchar(255);not null;default:''"`
	Birthday    time.Time  `gorm:"default:CURRENT_TIMESTAMP"`
	Phone       string     `gorm:"type:varchar(255);not null;unique_index"`
	Email       string     `gorm:"type:varchar(255);not null;unique;unique_index"`
	Address     string     `gorm:"type:varchar(255);not null;default:''"`
	Job         string     `gorm:"type:varchar(255);not null;default:''"`
	CountryCode string     `gorm:"type:varchar(255);not null;default:''"`
	Password    string     `gorm:"type:varchar(255);not null;default:''"`
	Role        value.Role `gorm:"type:tinyint(1);not null;default:1"`
	CreatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt   time.Time  `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func (user *SUser) BeforeCreate(tx *gorm.DB) (err error) {
	encryptedPwdData, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err == nil {
		user.Password = string(encryptedPwdData)
	}
	user.Username = strings.ToLower(html.EscapeString(strings.TrimSpace(user.Username)))

	return err
}
