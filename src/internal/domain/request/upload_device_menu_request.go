package request

type UploadDeviceMenuRequest struct {
	OrganizationID string                       `json:"organization_id" binding:"required"`
	DeviceID       string                       `json:"device_id" binding:"required"`
	Components     []CreateMenuComponentRequest `json:"components" binding:"required"`
}
