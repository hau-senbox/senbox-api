package request

type GetAllRoleClaimByRoleRequest struct {
	RoleId uint `json:"role_id" binding:"required"`
}

type GetRoleClaimByIdRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GetRoleClaimByNameRequest struct {
	ClaimName string `json:"claim_name" binding:"required"`
}
