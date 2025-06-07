package request

import "sen-global-api/internal/domain/entity/menu"

type CreateSuperAdminMenuRequest struct {
	Direction  menu.Direction               `json:"direction" binding:"required"`
	Components []CreateMenuComponentRequest `json:"components" binding:"required"`
}
