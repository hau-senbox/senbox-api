package entity

import (
	"github.com/google/uuid"
)

type SUserGuardians struct {
	UserID     uuid.UUID   `gorm:"column:user_id;primary_key"`
	User       SUserEntity `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE;"`
	GuardianId uuid.UUID   `gorm:"column:guardian_id;primary_key"`
	Guardian   SUserEntity `gorm:"foreignKey:GuardianId;references:id;constraint:OnDelete:CASCADE;"`
}
