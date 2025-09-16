package request

type UploadSectionMenuDepartmentRequest DepartmentSectionMenuItem

type DepartmentSectionMenuItem struct {
	Language           uint                         `json:"language" binding:"required"`
	DepartmentID       string                       `json:"department_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
