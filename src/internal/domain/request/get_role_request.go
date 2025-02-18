package request

type GetRoleByIdRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GetRoleByNameRequest struct {
	RoleName string `json:"role_name" binding:"required"`
}
