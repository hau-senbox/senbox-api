package request

type DeleteFunctionClaimPermissionRequest struct {
	ID uint `json:"id" binding:"required"`
}
