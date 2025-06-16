package request

type CreateTeacherFormApplicationRequest struct {
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
}
