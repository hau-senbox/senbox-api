package request

type GetDepartmentMenuOrganizationRequest struct {
	UserID         string `json:"user_id"`
	DeviceID       string `json:"device_id" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`
}
