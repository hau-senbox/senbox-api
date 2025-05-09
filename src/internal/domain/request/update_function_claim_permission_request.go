package request

type UpdateFunctionClaimPermissionRequest struct {
	PermissionId    uint   `json:"permission_id" binding:"required"`
	PermissionName  string `json:"permission_name" binding:"required"`
	FunctionClaimId int64  `json:"function_claim_id" binding:"required"`
}
