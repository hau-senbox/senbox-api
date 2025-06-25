package request

type GetUserDeviceByIDRequest struct {
	UserID   *string `json:"user_id"`
	DeviceID *string `json:"device_id"`
}
