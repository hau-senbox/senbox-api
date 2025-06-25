package request

type GetFunctionClaimPermissionByIDRequest struct {
	ID uint `json:"id" binding:"required"`
}

type GetFunctionClaimPermissionByNameRequest struct {
	PermissionName string `json:"permission_name" binding:"required"`
}
