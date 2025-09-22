package request

type UploadSectionDeviceMenuOrganizationRequest DevciceMenuOrganizationItem

type DevciceMenuOrganizationItem struct {
	LanguageID         uint                         `json:"language_id" binding:"required"`
	OrganizationID     string                       `json:"organization_id"`
	DeleteComponentIDs []string                     `json:"delete_component_ids"`
	Components         []CreateMenuComponentRequest `json:"components"`
}
