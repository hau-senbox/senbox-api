package request

type UploadOrgSettingPortalNewsRequest struct {
	OrganizationID     string `json:"organization_id"`
	IsPusblishedPortal bool   `json:"is_pusblished_portal"`
	MessagePortalNews  string `json:"message_portal_news" binding:"required"`
}
