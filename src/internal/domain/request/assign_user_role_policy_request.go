package request

type AssignUserRolePolicyRequest struct {
	User   string `json:"user_id" binding:"required"`
	Policy uint   `json:"policy_id" binding:"required"`
}
