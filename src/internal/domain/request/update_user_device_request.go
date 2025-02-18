package request

type UpdateUserDeviceRequest struct {
	UserId  string   `json:"user_id" binding:"required"`
	Devices []string `json:"devices" binding:"required"`
}
