package response

type OrgSettingDeviceNewsResponse struct {
	OrganizationID    string `json:"organization_id"`
	IsPublishedDevice bool   `json:"is_published_device"`
	MessageDeviceNews string `json:"message_devices_news"`
}

type OrgSettingPortalNewsResponse struct {
	OrganizationID    string `json:"organization_id"`
	IsPublishedPortal bool   `json:"is_published_portal"`
	MessagePortalNews string `json:"message_portal_news"`
}

type OrgSettingNewsResponse struct {
	OrganizationID    string `json:"organization_id"`
	IsPublishedDevice bool   `json:"is_published_device"`
	MessageDeviceNews string `json:"message_devices_news"`
	IsPublishedPortal bool   `json:"is_published_portal"`
	MessagePortalNews string `json:"message_portal_news"`
}
