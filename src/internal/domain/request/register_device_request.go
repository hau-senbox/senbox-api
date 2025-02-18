package request

type RegisterDeviceRequest struct {
	UserID     string `json:"user_id" binding:"required"`
	DeviceUUID string `json:"device_uuid" binding:"required"`
	InputMode  string `json:"input_mode" required:"true"`
	AppVersion string `json:"app_version" default:"" required:"true"`
}
