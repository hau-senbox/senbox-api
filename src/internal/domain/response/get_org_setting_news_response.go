package response

type OrgSettingDeviceNewsResponse struct {
	OrganizationID     string `json:"organization_id"`
	IsPusblishedDevice bool   `json:"is_pusblished_device"`
	MessageDeviceNews  string `json:"message_devices_news"`
}

type OrgSettingPortalNewsResponse struct {
	OrganizationID     string `json:"organization_id"`
	IsPusblishedPortal bool   `json:"is_pusblished_portal"`
	MessagePortalNews  string `json:"message_portal_news"`
}

type OrgSettingNewsResponse struct {
	OrganizationID     string `json:"organization_id"`
	IsPusblishedDevice bool   `json:"is_pusblished_device"`
	MessageDeviceNews  string `json:"message_devices_news"`
	IsPusblishedPortal bool   `json:"is_pusblished_portal"`
	MessagePortalNews  string `json:"message_portal_news"`
}
