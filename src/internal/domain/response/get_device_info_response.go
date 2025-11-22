package response

type GetDeviceInfoResponse struct {
	OrganizationID   string                  `json:"organization_id"`
	DeviceName       string                  `json:"device_name"`
	Components       []GetMenus4Web          `json:"components"`
	CreatedIndex     int                     `json:"created_index"`
	DeviceNickName   string                  `json:"device_nick_name"`
	CurrentAppValues *GetValuesAppResponse   `json:"current_app_values"`
	ValueHistories   []*GetValuesAppResponse `json:"value_histories"`
}

type GetPersonalDeviceInfoResponse struct {
	DeviceID       string                  `json:"device_id"`
	DeviceCode     string                  `json:"device_code"`
	Organizations  []OrganizationDevices   `json:"organizations"`
	Teachers       []TeacherResponse       `json:"teachers"`
	Staffs         []StaffResponse         `json:"staffs"`
	Students       []StudentResponse       `json:"students"`
	Parents        []ParentResponse        `json:"parents"`
	ValueHistories []*GetValuesAppResponse `json:"value_histories"`
}
