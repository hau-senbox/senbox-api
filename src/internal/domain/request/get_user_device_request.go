package request

type GetUserDeviceByIdRequest struct {
	UserId   *string `json:"user_id"`
	DeviceId *string `json:"device_id"`
}
