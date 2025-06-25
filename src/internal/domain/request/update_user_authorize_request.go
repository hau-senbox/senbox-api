package request

type UpdateUserAuthorizeRequest struct {
	UserID                    string `json:"user_id" binding:"required"`
	FunctionClaimID           int64  `json:"function_claim_id" binding:"required"`
	FunctionClaimPermissionID int64  `json:"function_claim_permission_id" binding:"required"`
}
