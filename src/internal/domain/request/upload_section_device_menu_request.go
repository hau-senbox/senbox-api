package request

type UploadSectionDeviceMenuOrganizationRequest DevciceMenuOrganizationItem

type DevciceMenuOrganizationItem struct {
	OrganizationID     string                       `json:"organization_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
