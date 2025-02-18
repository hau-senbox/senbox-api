package entity

import (
	"time"
)

type SRoleClaim struct {
	ID         int64     `gorm:"column:id;primary_key;AUTO_INCREMENT"`
	ClaimName  string    `gorm:"type:varchar(255);not null;"`
	ClaimValue string    `gorm:"type:varchar(255);not null;"`
	RoleId     int64     `gorm:"column:role_id;"`
	Role       SRole     `gorm:"foreignKey:RoleId;references:id;constraint:OnDelete:CASCADE"`
	CreatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt  time.Time `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
