package response

type ParentResponseBase struct {
	ParentID       string                   `json:"id"`
	ParentName     string                   `json:"name"`
	Avatar         string                   `json:"avatar,omitempty"`
	AvatarURL      string                   `json:"avatar_url,omitempty"`
	Menus          []GetMenus4Web           `json:"components"`
	CustomID       string                   `json:"custom_id"`
	LanguageConfig *LanguagesConfigResponse `json:"language_config"`
	Avatars        []Avatar                 `json:"avatars"`
	CreatedIndex   int                      `json:"created_index"`
	LogedDevices   []LoggedDevice           `json:"logged_devices"`
}

type GetParent4Gateway struct {
	ParentID       string `json:"id"`
	OrganizationID string `json:"organization_id"`
	ParentName     string `json:"name"`
	Avatar         Avatar `json:"avatar"`
	Code           string `json:"code"`
	UserID         string `json:"user_id"`
}
