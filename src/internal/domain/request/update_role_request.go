package request

type UpdateRoleRequest struct {
	RoleId   uint   `json:"role_id" binding:"required"`
	RoleName string `json:"role_name" binding:"required"`
}
