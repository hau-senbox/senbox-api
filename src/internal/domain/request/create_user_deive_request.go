package request

type RegisteringDeviceForUser struct {
	UserId   string `json:"user_id" binding:"required"`
	DeviceId string `json:"device_id" binding:"required"`
}
