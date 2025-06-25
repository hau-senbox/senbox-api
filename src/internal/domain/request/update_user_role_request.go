package request

type UpdateUserRoleRequest struct {
	UserID string   `json:"user_id" binding:"required"`
	Roles  []string `json:"roles" binding:"required"`
}
