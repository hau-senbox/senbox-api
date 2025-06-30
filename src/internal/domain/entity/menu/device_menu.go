package menu

import (
	"github.com/google/uuid"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
)

type DeviceMenu struct {
	OrganizationID uuid.UUID            `gorm:"column:organization_id;primary_key"`
	Organization   entity.SOrganization `gorm:"foreignKey:OrganizationID;references:id;constraint:OnDelete:CASCADE;"`
	ComponentID    uuid.UUID            `gorm:"column:component_id;primary_key"`
	Component      components.Component `gorm:"foreignKey:ComponentID;references:id;constraint:OnDelete:CASCADE;"`
	Order          int                  `gorm:"type:int;not null;default:0"`
}
