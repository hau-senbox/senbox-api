package request

type UploadSectionMenuDepartmentRequest DepartmentSectionMenuItem

type DepartmentSectionMenuItem struct {
	DepartmentID       string                       `json:"department_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
