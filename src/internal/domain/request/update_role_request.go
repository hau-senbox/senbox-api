package request

type UpdateRoleRequest struct {
	RoleID   uint   `json:"role_id" binding:"required"`
	RoleName string `json:"role_name" binding:"required"`
}
