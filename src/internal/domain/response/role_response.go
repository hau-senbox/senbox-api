package response

type RoleResponse struct {
	ID       int64  `json:"id"`
	RoleName string `json:"role_name"`
}

type RoleListResponseData struct {
	ID       int64  `json:"id"`
	RoleName string `json:"role_name"`
}
