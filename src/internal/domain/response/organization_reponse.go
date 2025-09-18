package response

import "sen-global-api/internal/domain/entity/components"

type OrgSettingResponse struct {
	ID                 string                `json:"id"`
	OrganizationID     string                `json:"organization_id" binding:"required"`
	OrganizationName   string                `json:"organization_name"`
	DeviceID           string                `json:"device_id" binding:"required"`
	IsViewMessageBox   bool                  `json:"is_view_message_box"`
	IsShowMessage      bool                  `json:"is_show_message"`
	MessageBox         string                `json:"message_box"`
	IsShowSpecialBtn   bool                  `json:"is_show_special_btn"`
	IsDeactiveApp      bool                  `json:"is_deactive_app"`
	MessageDeactiveApp string                `json:"message_deactive_app"`
	IsDeactiveTopMenu  bool                  `json:"is_deactive_top_menu"`
	MessageTopMenu     string                `json:"message_top_menu"`
	TopMenuPassword    string                `json:"top_menu_password"`
	Component          *components.Component `json:"component"`
}
