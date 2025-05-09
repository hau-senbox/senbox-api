package request

type CreateFunctionClaimRequest struct {
	FunctionName string `json:"function_name" binding:"required"`
}

type CreateFunctionClaimsRequest struct {
	FunctionClaims []CreateFunctionClaimRequest `json:"function_claims" binding:"required"`
}
