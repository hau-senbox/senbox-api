package request

type RegisterDeviceRequest struct {
	DeviceUUID string `json:"device_uuid" binding:"required"`
	InputMode  string `json:"input_mode" required:"true"`
	AppVersion string `json:"app_version" default:"" required:"true"`
}
