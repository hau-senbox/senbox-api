package request

type UserLoginFromDeviceReqest struct {
	Type       string `json:"type"`
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	DeviceUUID string `json:"device_uuid" binding:"required"`
}

type UserLoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserLogoutReqeust struct {
	DeviceID string `json:"device_id" binding:"required"`
}
