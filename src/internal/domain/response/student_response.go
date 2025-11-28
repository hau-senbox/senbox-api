package response

import gw_response "sen-global-api/pkg/consulapi/gateway/dto/response"

type StudentResponseBase struct {
	StudentID      string                          `json:"id"`
	Code           string                          `json:"code"`
	StudentName    string                          `json:"name"`
	Avatar         string                          `json:"avatar,omitempty"`
	AvatarURL      string                          `json:"avatar_url"`
	QrFormProfile  string                          `json:"qr_form,omitempty"`
	Menus          []GetMenus4Web                  `json:"components"`
	CustomID       string                          `json:"custom_id"`
	StudentBlock   *StudentBlockSettingResponse    `json:"student_block"`
	LanguageConfig *LanguagesConfigResponse        `json:"language_config"`
	Avatars        []Avatar                        `json:"avatars"`
	CreatedIndex   int                             `json:"created_index"`
	LogedDevices   []LoggedDevice                  `json:"logged_devices"`
	Information    *gw_response.StudentInformation `json:"information"`
}

type LoggedDevice struct {
	DeviceID   string `json:"device_id"`
	DeviceCode string `json:"device_code"`
}

type GetStudent4Gateway struct {
	StudentID      string `json:"id"`
	OrganizationID string `json:"organization_id"`
	StudentName    string `json:"name"`
	Avatar         Avatar `json:"avatar"`
	Code           string `json:"code"`
}
