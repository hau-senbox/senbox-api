package request

type UpdateRoleRequest struct {
	ID          uint   `json:"id" binding:"required"`
	RoleName    string `json:"role_name" binding:"required"`
	Description string `json:"desciption" binding:"required"`
}
