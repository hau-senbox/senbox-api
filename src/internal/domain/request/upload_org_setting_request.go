package request

type UploadOrgSettingRequest struct {
	OrganizationID    string                     `json:"organization_id" binding:"required,uuid"`
	DeviceID          string                     `json:"device_id" binding:"required"`
	IsViewMessage     bool                       `json:"is_view_message"`
	IsShowOrgNews     bool                       `json:"is_show_org_news"`
	IsDeactiveTopMenu bool                       `json:"is_deactive_top_menu"`
	IsShowSpecialBtn  bool                       `json:"is_show_special_btn"`
	MessageBox        string                     `json:"message_box" binding:"max=500"`
	MessageTopMenu    string                     `json:"message_top_menu" binding:"max=500"`
	TopMenuPassword   string                     `json:"top_menu_password"`
	Component         CreateMenuComponentRequest `json:"component" binding:"required"`
}
