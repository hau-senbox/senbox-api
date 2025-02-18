package request

type UpdateRolePolicyRequest struct {
	ID          uint    `json:"id" binding:"required"`
	PolicyName  string  `json:"policy_name" binding:"required"`
	Description string  `json:"desciption" binding:"required"`
	Roles       *[]uint `json:"roles"`
	RoleClaims  *[]uint `json:"role_claims"`
}
