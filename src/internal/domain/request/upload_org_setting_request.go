package request

type UploadOrgSettingRequest struct {
	OrganizationID     string                      `json:"organization_id"`
	DeviceID           string                      `json:"device_id"`
	IsViewMessageBox   *bool                       `json:"is_view_message_box"`
	IsShowMessage      *bool                       `json:"is_show_message"`
	MessageBox         *string                     `json:"message_box"`
	IsShowSpecialBtn   *bool                       `json:"is_show_special_btn"`
	IsDeactiveApp      *bool                       `json:"is_deactive_app"`
	MessageDeactiveApp *string                     `json:"message_deactive_app"`
	IsDeactiveTopMenu  *bool                       `json:"is_deactive_top_menu"`
	MessageTopMenu     *string                     `json:"message_top_menu"`
	TopMenuPassword    *string                     `json:"top_menu_password"`
	Component          UploadOrgSettingMenuRequest `json:"component"`
}

type UploadOrgSettingMenuRequest struct {
	ID     string                    `json:"id"`
	Name   string                    `json:"name"`
	Type   string                    `json:"type"`
	Key    string                    `json:"key" default:""`
	Value  UploadOrgSettingMenuValue `json:"value"`
	Order  int                       `json:"order"`
	IsShow bool                      `json:"is_show"`
}

type UploadOrgSettingMenuValue struct {
	Icon         string `json:"icon"`
	Visible      bool   `json:"visible"`
	Color        string `json:"color"`
	URL          string `json:"url"`
	FormQr       string `json:"form_qr"`
	ShowedTop    bool   `json:"showed_top"`
	ShowedBottom bool   `json:"showed_bottom"`
	Note         string `json:"note"`
}
