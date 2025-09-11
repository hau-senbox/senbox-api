package response

type GetDeviceInfoResponse struct {
	OrganizationID   string                       `json:"organization_id"`
	DeviceName       string                       `json:"device_name"`
	Components       []ComponentResponse          `json:"components"`
	CreatedIndex     int                          `json:"created_index"`
	DeviceNickName   string                       `json:"device_nick_name"`
	ValuesAppCurrent *GetValuesAppCurrentResponse `json:"values_app_current"`
}
