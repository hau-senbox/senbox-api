package request

type GetRoleByIDRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GetRoleByNameRequest struct {
	RoleName string `json:"role" binding:"required"`
}
