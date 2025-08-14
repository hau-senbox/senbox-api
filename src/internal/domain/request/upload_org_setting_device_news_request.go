package request

type UploadOrgSettingDeviceNewsRequest struct {
	OrganizationID     string `json:"organization_id"`
	IsPusblishedDevice bool   `json:"is_pusblished_device"`
	MessageDeviceNews  string `json:"message_devices_news" binding:"required"`
}
