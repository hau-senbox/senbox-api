package request

type AssignUserRoleRequest struct {
	User string `json:"user_id" binding:"required"`
	Role uint   `json:"role_id" binding:"required"`
}
