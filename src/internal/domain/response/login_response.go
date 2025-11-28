package response

import "time"

type LoginResponseData struct {
	UserID            string             `json:"user_id"`
	Username          string             `json:"username"`
	IsSuperAdmin      bool               `json:"is_super_admin"`
	Organizations     []string           `json:"organizations"`
	Token             string             `json:"token"`
	Expired           time.Time          `json:"expired"`
	OrganizationAdmin *OrganizationAdmin `json:"organization_admin"`
	RedirectUrl       string             `json:"redirect_url"`
	AllowedRouters    []AllowedRouters   `json:"allowed_routers"`
}

type LoginResponse struct {
	Data LoginResponseData `json:"data"`
}
