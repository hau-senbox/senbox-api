package response

type RoleResponse struct {
	ID          int64  `json:"id"`
	RoleName    string `json:"role_name"`
	Description string `json:"description"`
}

type RoleListResponseData struct {
	ID       int64  `json:"id"`
	RoleName string `json:"role_name"`
}

type RoleListResponse struct {
	Data []RoleListResponseData `json:"data"`
}
