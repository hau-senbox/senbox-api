package request

type CreateDeviceMenuRequest struct {
	DeviceID       string                       `json:"device_id" binding:"required"`
	OrganizationID string                       `json:"organization_id" binding:"required"`
	Components     []CreateMenuComponentRequest `json:"components" binding:"required"`
}
