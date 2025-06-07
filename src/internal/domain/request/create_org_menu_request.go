package request

import "sen-global-api/internal/domain/entity/menu"

type CreateOrgMenuRequest struct {
	OrganizationID int64                        `json:"organization_id" binding:"required"`
	Direction      menu.Direction               `json:"direction" binding:"required"`
	Components     []CreateMenuComponentRequest `json:"components" binding:"required"`
}
