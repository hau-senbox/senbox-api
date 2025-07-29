package response

type TeacherFormApplicationResponse struct {
	ID               string `json:"id"`
	TeacherName      string `json:"teacher_name"`
	Status           string `json:"status"`
	ApprovedAt       string `json:"approved_at"`
	CreatedAt        string `json:"created_at"`
	UserID           string `json:"user_id"`
	OrganizationID   string `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
}
