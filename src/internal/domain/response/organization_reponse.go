package response

import "sen-global-api/internal/domain/entity/components"

type OrgSettingResponse struct {
	ID                string                `json:"id"`
	OrganizationID    string                `json:"organization_id"`
	DeviceID          string                `json:"device_id"`
	IsViewMessage     bool                  `json:"is_view_message"`
	IsShowOrgNews     bool                  `json:"is_show_org_news"`
	IsDeactiveTopMenu bool                  `json:"is_deactive_top_menu"`
	IsShowSpecialBtn  bool                  `json:"is_show_special_btn"`
	MessageBox        string                `json:"message_box"`
	MessageTopMenu    string                `json:"message_top_menu"`
	TopMenuPasswod    string                `json:"top_menu_password"`
	Component         *components.Component `json:"component"`
}
