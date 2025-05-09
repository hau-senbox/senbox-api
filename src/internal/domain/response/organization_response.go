package response

type OrganizationResponse struct {
	ID          int64  `json:"id"`
	OrganizationName string `json:"organization_name"`
	Address     string `json:"address"`
	Description string `json:"description"`
}
