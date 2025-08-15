package response

type GetDeviceInfoResponse struct {
	OrganizationID string              `json:"organization_id"`
	DeviceName     string              `json:"device_name"`
	Components     []ComponentResponse `json:"components"`
}
