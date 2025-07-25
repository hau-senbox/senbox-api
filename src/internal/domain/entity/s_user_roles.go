package entity

import (
	"github.com/google/uuid"
)

type SUserRoles struct {
	UserID uuid.UUID   `gorm:"column:user_id;primary_key"`
	User   SUserEntity `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE;"`
	RoleID int64       `gorm:"column:role_id;primary_key"`
	Role   SRole       `gorm:"foreignKey:RoleID;references:id;constraint:OnDelete:CASCADE"`
}
