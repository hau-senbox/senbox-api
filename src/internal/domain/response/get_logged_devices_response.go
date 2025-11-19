package response

type GetLoggedDevicesResponse struct {
	DeviceID            string                `json:"device_id"`
	DeviceCode          string                `json:"device_code"`
	OrganizationDevices []OrganizationDevices `json:"organization_devices"`
}

type OrganizationDevices struct {
	OrganizationID         string `json:"organization_id"`
	OrganizationName       string `json:"organization_name"`
	OrganizationDeviceCode string `json:"organization_device_code"`
}
