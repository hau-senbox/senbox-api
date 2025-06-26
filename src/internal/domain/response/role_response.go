package response

type RoleResponse struct {
	ID       int64  `json:"id"`
	RoleName string `json:"role"`
}

type RoleListResponseData struct {
	ID       int64  `json:"id"`
	RoleName string `json:"role"`
}
