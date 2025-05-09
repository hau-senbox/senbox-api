package request

type DeleteUserAuthorizeRequest struct {
	UserId          string `json:"user_id" binding:"required"`
	FunctionClaimId int64  `json:"function_claim_id" binding:"required"`
}
