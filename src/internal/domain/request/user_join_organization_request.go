package request

type UserJoinOrganizationRequest struct {
	UserId         string `json:"user_id" binding:"required"`
	OrganizationId uint   `json:"organization_id" binding:"required"`
	Password       string `json:"password" binding:"required"`
}
