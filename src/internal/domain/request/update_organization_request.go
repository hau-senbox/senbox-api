package request

type UpdateOrganizationRequest struct {
	ID          int64  `json:"id" binding:"required"`
	OrganizationName string `json:"organization_name" binding:"required"`
	Address     string `json:"address" binding:"required"`
	Description string `json:"description" binding:"required"`
}
