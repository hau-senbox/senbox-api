package response

import "time"

type SwitchToOrganizationResponse struct {
	Token   string               `json:"token"`
	Expired time.Time            `json:"expired"`
	User    UserEntityResponseV2 `json:"user"`
}
