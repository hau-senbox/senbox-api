package request

type UploadOrgSettingDeviceNewsRequest struct {
	OrganizationID string `json:"organization_id"`
	IsNewsConfig   bool   `json:"is_news_config"`
	IsPusblished   bool   `json:"is_pusblished"`
	MessageNews    string `json:"message_news" binding:"required"`
}
