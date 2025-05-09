package request

type CreateFunctionClaimPermissionRequest struct {
	PermissionName  string `json:"permission_name" binding:"required"`
	FunctionClaimId int64  `json:"function_claim_id" binding:"required"`
}
