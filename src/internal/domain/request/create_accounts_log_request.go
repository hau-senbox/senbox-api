package request

type CreateAccountsLogRequest struct {
	Type           string `json:"type"`
	UserID         string `json:"user_id"`
	DeviceID       string `json:"device_id"`
	OrganizationID string `json:"organization_id"`
}
