package request

type UpdateUserAuthorizeRequest struct {
	UserId                    string `json:"user_id" binding:"required"`
	FunctionClaimId           int64  `json:"function_claim_id" binding:"required"`
	FunctionClaimPermissionId int64  `json:"function_claim_permission_id" binding:"required"`
}
