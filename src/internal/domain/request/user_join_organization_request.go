package request

type UserJoinOrganizationRequest struct {
	UserID         string `json:"user_id" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`
	Password       string `json:"password" default:""`
}
