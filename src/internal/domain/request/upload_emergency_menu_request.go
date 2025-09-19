package request

type UploadEmergencyMenuRequest EmergencyMenuMenuItem

type EmergencyMenuMenuItem struct {
	LanguageID         uint                         `json:"language_id" binding:"required"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
