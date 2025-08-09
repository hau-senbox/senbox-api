package response

type GetDeviceMenuResponse struct {
	DeviceID   string              `json:"device_id"`
	DeviceName string              `json:"device_name"`
	Components []ComponentResponse `json:"components"`
}
