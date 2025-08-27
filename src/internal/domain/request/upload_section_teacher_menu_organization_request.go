package request

type UploadSectionMenuTeacherOrganizationRequest TeacherOrganizationSectionMenuItem

type TeacherOrganizationSectionMenuItem struct {
	TeacherID          string                       `json:"teacher_id"`
	OrganizationID     string                       `json:"organization_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
