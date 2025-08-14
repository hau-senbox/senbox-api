package response

type GetOrgSettingNewsResponse struct {
	OrganizationID string `json:"organization_id"`
	IsPusblished   bool   `json:"is_pusblished"`
	MessageNews    string `json:"message_news" binding:"required"`
}
