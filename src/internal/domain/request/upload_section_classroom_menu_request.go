package request

type UploadSectionMenuClassroomRequest ClassroomSectionMenuItem

type ClassroomSectionMenuItem struct {
	ClassroomID        string                       `json:"classroom_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
