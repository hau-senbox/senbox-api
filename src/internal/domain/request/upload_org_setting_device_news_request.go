package request

type UploadOrgSettingDeviceNewsRequest struct {
	OrganizationID    string `json:"organization_id"`
	IsPublishedDevice bool   `json:"is_published_device"`
	MessageDeviceNews string `json:"message_device_news" binding:"required"`
}
