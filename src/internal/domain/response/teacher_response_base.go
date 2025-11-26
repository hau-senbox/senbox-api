package response

type TeacherResponseBase struct {
	TeacherID      string                   `json:"id"`
	UserID         string                   `json:"user_id"`
	TeacherName    string                   `json:"name"`
	Avatar         string                   `json:"avatar,omitempty"`
	AvatarURL      string                   `json:"avatar_url,omitempty"`
	QrFormProfile  string                   `json:"qr_form,omitempty"`
	Menus          []GetMenus4Web           `json:"components"`
	IsUserBlock    bool                     `json:"is_block"`
	LanguageConfig *LanguagesConfigResponse `json:"language_config"`
	Avatars        []Avatar                 `json:"avatars"`
	CreatedIndex   int                      `json:"created_index"`
	LogedDevices   []LoggedDevice           `json:"logged_devices"`
}

type GetTeacher4Gateway struct {
	TeacherID      string `json:"id"`
	OrganizationID string `json:"organization_id"`
	TeacherName    string `json:"name"`
	Avatar         Avatar `json:"avatar"`
	Code           string `json:"code"`
	UserID         string `json:"user_id"`
}
