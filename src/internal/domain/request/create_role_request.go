package request

type CreateRoleRequest struct {
	RoleName string `json:"role" binding:"required"`
}
