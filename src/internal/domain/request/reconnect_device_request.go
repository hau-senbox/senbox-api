package request

type ReconnectDeviceRequest struct {
	DeviceId string `json:"device_id" binding:"required"`
}
