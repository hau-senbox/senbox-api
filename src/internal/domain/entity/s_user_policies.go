package entity

import (
	"github.com/google/uuid"
)

type SUserPolicies struct {
	UserId     uuid.UUID   `gorm:"column:user_id;primary_key"`
	User       SUserEntity `gorm:"foreignKey:UserId;references:id;constraint:OnDelete:CASCADE;"`
	PolicyId   int64       `gorm:"column:policy_id;primary_key"`
	RolePolicy SRolePolicy `gorm:"foreignKey:PolicyId;references:id;constraint:OnDelete:CASCADE;"`
}
