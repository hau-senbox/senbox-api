package menu

import (
	"github.com/google/uuid"
	"sen-global-api/internal/domain/entity/components"
)

type SuperAdminMenu struct {
	Direction   Direction            `gorm:"column:direction;primary_key"`
	ComponentID uuid.UUID            `gorm:"column:component_id;primary_key"`
	Component   components.Component `gorm:"foreignKey:ComponentID;references:id;constraint:OnDelete:CASCADE;"`
	Order       int                  `gorm:"type:int;not null;default:0"`
}
