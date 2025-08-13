package request

import "github.com/google/uuid"

type UploadOrgSettingRequest struct {
	OrganizationID     string                      `json:"organization_id"`
	DeviceID           string                      `json:"device_id"`
	IsViewMessageBox   bool                        `json:"is_view_message_box"`
	IsShowMessage      bool                        `json:"is_show_message"`
	MessageBox         string                      `json:"message_box"`
	IsShowSpecialBtn   bool                        `json:"is_show_special_btn"`
	IsDeactiveApp      bool                        `json:"is_deactive_app"`
	MessageDeactiveApp string                      `json:"message_deactive_app"`
	IsDeactiveTopMenu  bool                        `json:"is_deactive_top_menu"`
	MessageTopMenu     string                      `json:"message_top_menu"`
	TopMenuPassword    string                      `json:"top_menu_password"`
	Component          UploadOrgSettingMenuRequest `json:"component"`
}

type UploadOrgSettingMenuRequest struct {
	ID     *uuid.UUID `json:"id"`
	Name   string     `json:"name"`
	Type   string     `json:"type"`
	Key    string     `json:"key" default:""`
	Value  string     `json:"value"`
	Order  int        `json:"order"`
	IsShow bool       `json:"is_show"`
}
