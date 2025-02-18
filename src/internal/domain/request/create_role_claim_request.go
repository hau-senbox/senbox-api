package request

type CreateRoleClaimRequest struct {
	ClaimName  string `json:"claim_name" binding:"required"`
	ClaimValue string `json:"claim_value" binding:"required"`
	RoleId     uint   `json:"role_id" binding:"required"`
}

type CreateRoleClaimsRequest struct {
	RoleClaims []CreateRoleClaimRequest `json:"role_claims" binding:"required"`
}
