package request

type CreateRoleRequest struct {
	RoleName    string `json:"role_name" binding:"required"`
	Description string `json:"desciption" default:"" binding:"required"`
}
