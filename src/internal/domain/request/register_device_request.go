package request

type RegisterDeviceUser struct {
	Fullname string `json:"fullname" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Phone    string `json:"phone" binding:"required"`
}

type RegisterDeviceRequest struct {
	Primary           RegisterDeviceUser `json:"primary" binding:"required"`
	Secondary         RegisterDeviceUser `json:"secondary" binding:"required"`
	Tertiary          RegisterDeviceUser `json:"tertiary" binding:"required"`
	DeviceUUID        string             `json:"device_uuid" binding:"required"`
	ProfilePictureUrl string             `json:"profile_picture_url"`
	InputMode         string             `json:"input_mode" required:"true"`
	AppVersion        string             `json:"app_version" default:"" required:"true"`
}
