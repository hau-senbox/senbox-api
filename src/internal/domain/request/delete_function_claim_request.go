package request

type DeleteFunctionClaimRequest struct {
	ID uint `json:"id" binding:"required"`
}
