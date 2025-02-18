package request

type UpdateRoleClaimRequest struct {
	ID         uint   `json:"id" binding:"required"`
	ClaimName  string `json:"claim_name" binding:"required"`
	ClaimValue string `json:"claim_value" binding:"required"`
}
