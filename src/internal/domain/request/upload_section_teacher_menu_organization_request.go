package request

type UploadSectionMenuTeacherOrganizationRequest TeacherOrganizationSectionMenuItem

type TeacherOrganizationSectionMenuItem struct {
	LanguageID         uint                         `json:"language_id" binding:"required"`
	TeacherID          string                       `json:"teacher_id"`
	OrganizationID     string                       `json:"organization_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}

type UploadSectionMenuStudentOrganizationRequest StudentOrganizationSectionMenuItem

type StudentOrganizationSectionMenuItem struct {
	LanguageID         uint                         `json:"language_id" binding:"required"`
	StudentID          string                       `json:"student_id"`
	OrganizationID     string                       `json:"organization_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
