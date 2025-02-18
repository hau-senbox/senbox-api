package request

type CreateRolePolicyRequest struct {
	PolicyName  string  `json:"policy_name" binding:"required"`
	Description string  `json:"desciption" default:"" binding:"required"`
	Roles       *[]uint `json:"roles"`
	RoleClaims  *[]uint `json:"role_claims"`
}
