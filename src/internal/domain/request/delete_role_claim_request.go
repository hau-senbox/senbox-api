package request

type DeleteRoleClaimRequest struct {
	ID uint `json:"id" binding:"required"`
}
