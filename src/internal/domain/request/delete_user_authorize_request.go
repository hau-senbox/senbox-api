package request

type DeleteUserAuthorizeRequest struct {
	UserID          string `json:"user_id" binding:"required"`
	FunctionClaimID int64  `json:"function_claim_id" binding:"required"`
}
