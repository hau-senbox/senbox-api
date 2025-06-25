package request

type UpdateUserDeviceRequest struct {
	UserID  string   `json:"user_id" binding:"required"`
	Devices []string `json:"devices" binding:"required"`
}
