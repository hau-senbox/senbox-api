package request

type GetFunctionClaimByIDRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GetFunctionClaimByNameRequest struct {
	FunctionName string `json:"function_name" binding:"required"`
}
