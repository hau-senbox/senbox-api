package request

type UserJoinOrganizationRequest struct {
	UserId         string `json:"user_id" binding:"required"`
	OrganizationId int64  `json:"organization_id" binding:"required"`
	Password       string `json:"password" default:""`
}
