package response

type OrganizationResponse struct {
	ID               string                      `json:"id"`
	OrganizationName string                      `json:"organization_name"`
	Address          string                      `json:"address"`
	Description      string                      `json:"description"`
	Managers         []GetOrgManagerInfoResponse `json:"managers"`
}
