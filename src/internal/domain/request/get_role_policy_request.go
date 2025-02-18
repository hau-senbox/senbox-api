package request

type GetRolePolicyByIdRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GetRolePolicyByNameRequest struct {
	PolicyName string `json:"policy_name" binding:"required"`
}
