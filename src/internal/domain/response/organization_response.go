package response

type OrganizationResponse struct {
	ID               string                      `json:"id"`
	Code             string                      `json:"code"`
	OrganizationName string                      `json:"organization_name"`
	Avatar           string                      `json:"avatar"`
	AvatarURL        string                      `json:"avatar_url"`
	Address          string                      `json:"address"`
	Description      string                      `json:"description"`
	Managers         []GetOrgManagerInfoResponse `json:"managers"`
}
