package request

type CreateOrgFormApplicationRequest struct {
	OrganizationName   string `json:"organization_name" binding:"required"`
	ApplicationContent string `json:"application_content" binding:"required"`
	UserID             string `json:"user_id"`
}
