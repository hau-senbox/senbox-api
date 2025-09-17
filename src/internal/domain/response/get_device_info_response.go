package response

type GetDeviceInfoResponse struct {
	OrganizationID   string                       `json:"organization_id"`
	DeviceName       string                       `json:"device_name"`
	Components       []GetMenus4Web               `json:"components"`
	CreatedIndex     int                          `json:"created_index"`
	DeviceNickName   string                       `json:"device_nick_name"`
	CurrentAppValues *GetValuesAppCurrentResponse `json:"current_app_values"`
}
