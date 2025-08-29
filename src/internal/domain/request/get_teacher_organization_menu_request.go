package request

type GetTeacherOrganizationMenuRequest struct {
	UserID         string `json:"user_id" binding:"required"`
	OrganizationID string `json:"organization_id" binding:"required"`
	DeviceID       string `json:"device_id" binding:"required"`
}
