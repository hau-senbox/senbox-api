package request

type UploadEmergencyMenuRequest EmergencyMenuMenuItem

type EmergencyMenuMenuItem struct {
	Language           uint                         `json:"language"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
