package request

type UpdateUserOrgInfoRequest struct {
	UserID         string `json:"user_id" binding:"required"`
	OrganizationID int64  `json:"organization_id" binding:"required"`
	UserNickName   string `json:"user_nick_name" binding:"required"`
	IsManager      bool   `json:"is_manager" binding:"required"`
}
