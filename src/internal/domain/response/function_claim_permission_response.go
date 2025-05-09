package response

type FunctionClaimPermissionResponse struct {
	ID             int64  `json:"id"`
	PermissionName string `json:"permission_name"`
}

type FunctionClaimPermissionListResponseData struct {
	ID             int64  `json:"id"`
	PermissionName string `json:"permission_name"`
}
