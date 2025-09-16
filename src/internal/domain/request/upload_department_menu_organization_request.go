package request

type UploadDepartmentMenuOrganizationRequest DepartmentSectionMenuOrganizationItem

type DepartmentSectionMenuOrganizationItem struct {
	Language           uint                         `json:"language" binding:"required"`
	DepartmentID       string                       `json:"department_id"`
	OrganizationID     string                       `json:"organization_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
