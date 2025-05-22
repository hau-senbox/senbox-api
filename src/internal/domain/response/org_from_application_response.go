package response

type OrgFormApplicationResponse struct {
	ID                 int64  `json:"id"`
	OrganizationName   string `json:"organization_name"`
	ApplicationContent string `json:"application_content"`
	Status             string `json:"status"`
	ApprovedAt         string `json:"approved_at"`
	CreatedAt          string `json:"created_at"`
	UserId             string `json:"user_id"`
}
