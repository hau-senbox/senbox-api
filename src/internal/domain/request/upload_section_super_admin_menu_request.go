package request

import "sen-global-api/internal/domain/entity/menu"

type UploadSectionSuperAdminMenuRequest SuperAdminSectionMenuItem

type SuperAdminSectionMenuItem struct {
	LanguageID         uint                         `json:"language_id" binding:"required"`
	Direction          menu.Direction               `json:"direction"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
