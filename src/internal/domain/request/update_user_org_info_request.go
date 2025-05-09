package request

type UpdateUserOrgInfoRequest struct {
	UserId         string `json:"user_id" binding:"required"`
	OrganizationId string `json:"organization_id" binding:"required"`
	UserNickName   string `json:"user_nick_name" binding:"required"`
	IsManager      bool   `json:"is_manager" binding:"required"`
}
