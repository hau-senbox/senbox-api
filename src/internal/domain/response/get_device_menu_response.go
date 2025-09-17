package response

type GetDeviceMenuResponse struct {
	DeviceID   string         `json:"device_id"`
	DeviceName string         `json:"device_name"`
	Components []GetMenus4Web `json:"components"`
}
