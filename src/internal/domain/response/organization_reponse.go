package response

import "sen-global-api/internal/domain/entity/components"

type OrgSettingResponse struct {
	ID                string                `json:"id"`
	OrganizationID    string                `json:"organization_id"`
	IsViewMessage     bool                  `json:"is_view_message"`
	IsShowOrgNews     bool                  `json:"is_show_org_news"`
	IsDeactiveTopMenu bool                  `json:"is_deactive_top_menu"`
	IsShowSpecialBtn  bool                  `json:"is_show_special_btn"`
	MessageBox        string                `json:"message_box"`
	MessageTopMenu    string                `json:"message_top_menu"`
	Component         *components.Component `json:"component"`
	CreatedAt         string                `json:"created_at"`
	UpdatedAt         string                `json:"updated_at"`
}
