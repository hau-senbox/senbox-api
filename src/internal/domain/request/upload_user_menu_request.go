package request

type UploadUserMenuRequest struct {
	OrganizationID string                       `json:"organization_id" binding:"required"`
	UserID         string                       `json:"user_id" binding:"required"`
	Components     []CreateMenuComponentRequest `json:"components" binding:"required"`
}
