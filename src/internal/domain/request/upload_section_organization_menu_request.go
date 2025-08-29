package request

import "sen-global-api/internal/domain/entity/menu"

type UploadSectionOrganizationAdminMenuRequest OrganizationAdminSectionMenuItem

type OrganizationAdminSectionMenuItem struct {
	Direction          menu.Direction               `json:"direction"`
	OrganizationID     string                       `json:"organization_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
