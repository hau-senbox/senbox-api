package request

type UpdateUserRolePolicyRequest struct {
	UserId   string `json:"user_id" binding:"required"`
	Policies []uint `json:"policies" binding:"required"`
}
