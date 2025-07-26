package menu

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"

	"github.com/google/uuid"
)

type UserMenu struct {
	UserID      uuid.UUID            `gorm:"column:user_id;primary_key"`
	User        entity.SUserEntity   `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE;"`
	ComponentID uuid.UUID            `gorm:"column:component_id;primary_key"`
	Component   components.Component `gorm:"foreignKey:ComponentID;references:id;constraint:OnDelete:CASCADE;"`
	Order       int                  `gorm:"type:int;not null;default:0"`
}
