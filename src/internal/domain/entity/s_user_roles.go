package entity

import (
	"github.com/google/uuid"
)

type SUserRoles struct {
	UserId uuid.UUID   `gorm:"column:user_id;primary_key"`
	User   SUserEntity `gorm:"foreignKey:UserId;references:id;constraint:OnDelete:CASCADE;"`
	RoleId int64       `gorm:"column:role_id;primary_key"`
	Role   SRole       `gorm:"foreignKey:RoleId;references:id;constraint:OnDelete:CASCADE"`
}
