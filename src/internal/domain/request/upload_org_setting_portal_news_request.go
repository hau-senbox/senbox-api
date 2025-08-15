package request

type UploadOrgSettingPortalNewsRequest struct {
	OrganizationID    string `json:"organization_id"`
	IsPublishedPortal bool   `json:"is_published_portal"`
	MessagePortalNews string `json:"message_portal_news" binding:"required"`
}
