package request

type CreateTeacherFormApplicationRequest struct {
	UserID         string `json:"user_id"`
	OrganizationID int64  `json:"organization_id"`
}
