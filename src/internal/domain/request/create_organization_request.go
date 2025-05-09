package request

type CreateOrganizationRequest struct {
	OrganizationName string `json:"organization_name" binding:"required"`
	Password         string `json:"password" binding:"required"`
	Address          string `json:"address" default:""`
	Description      string `json:"description" default:""`
}
