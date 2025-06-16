package request

type UploadOrgMenuRequest struct {
	OrganizationID string                       `json:"organization_id" binding:"required"`
	Top            []CreateMenuComponentRequest `json:"top" binding:"required"`
	Bottom         []CreateMenuComponentRequest `json:"bottom" binding:"required"`
}
