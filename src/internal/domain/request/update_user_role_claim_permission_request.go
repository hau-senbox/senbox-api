package request

type UpdateUserRoleClaimPermissionRequest struct {
	UserId               string  `json:"user_id" binding:"required"`
	Roles                []int64 `json:"roles" binding:"required"`
	RoleCalims           []int64 `json:"role_claims" binding:"required"`
	RoleCalimPermissions []int64 `json:"role_claim_permissions" binding:"required"`
}
