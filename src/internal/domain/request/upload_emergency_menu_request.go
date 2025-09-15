package request

type UploadEmergencyMenuRequest EmergencyMenuMenuItem

type EmergencyMenuMenuItem struct {
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
