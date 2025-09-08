package request

type UploadDepartmentMenuOrganizationRequest DepartmentSectionMenuOrganizationItem

type DepartmentSectionMenuOrganizationItem struct {
	DepartmentID       string                       `json:"department_id"`
	OrganizationID     string                       `json:"organization_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
