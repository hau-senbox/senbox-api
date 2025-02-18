package request

type DeleteRoleRequest struct {
	ID uint `json:"id" binding:"required"`
}
