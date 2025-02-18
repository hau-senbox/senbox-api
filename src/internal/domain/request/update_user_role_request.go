package request

type UpdateUserRoleRequest struct {
	UserId string `json:"user_id" binding:"required"`
	Roles  []uint `json:"roles" binding:"required"`
}
