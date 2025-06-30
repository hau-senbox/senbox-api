package request

type CreateDeviceMenuRequest struct {
	OrganizationID string                       `json:"organization_id" binding:"required"`
	Components     []CreateMenuComponentRequest `json:"components" binding:"required"`
}
