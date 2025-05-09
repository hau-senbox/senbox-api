package request

type UpdateFunctionClaimRequest struct {
	FunctionClaimId uint   `json:"function_claim_id" binding:"required"`
	FunctionName    string `json:"function_name" binding:"required"`
}
