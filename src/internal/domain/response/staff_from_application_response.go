package response

type StaffFormApplicationResponse struct {
	ID               string `json:"id"`
	Status           string `json:"status"`
	ApprovedAt       string `json:"approved_at"`
	CreatedAt        string `json:"created_at"`
	UserID           string `json:"user_id"`
	StaffName        string `json:"staff_name"`
	OrganizationID   string `json:"organization_id"`
	OrganizationName string `json:"organization_name"`
}
