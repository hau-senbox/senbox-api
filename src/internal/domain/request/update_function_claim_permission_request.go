package request

type UpdateFunctionClaimPermissionRequest struct {
	PermissionID    uint   `json:"permission_id" binding:"required"`
	PermissionName  string `json:"permission_name" binding:"required"`
	FunctionClaimID int64  `json:"function_claim_id" binding:"required"`
}
