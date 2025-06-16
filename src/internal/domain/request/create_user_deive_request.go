package request

type RegisteringDeviceForUser struct {
	UserID   string `json:"user_id" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}

type RegisteringDeviceForOrg struct {
	OrgID    string `json:"organization_id" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}
