package response

type StudentResponseBase struct {
	StudentID      string                       `json:"id"`
	StudentName    string                       `json:"name"`
	Avatar         string                       `json:"avatar,omitempty"`
	AvatarURL      string                       `json:"avatar_url,omitempty"`
	QrFormProfile  string                       `json:"qr_form,omitempty"`
	Menus          []ComponentResponse          `json:"components"`
	CustomID       string                       `json:"custom_id"`
	StudentBlock   *StudentBlockSettingResponse `json:"student_block"`
	LanguageConfig *LanguagesConfigResponse     `json:"language_config"`
	Avatars        []Avatar                     `json:"avatars"`
}

type GetStudent4Gateway struct {
	StudentID      string `json:"id"`
	OrganizationID string `json:"organization_id"`
	StudentName    string `json:"name"`
	Avatar         Avatar `json:"avatar"`
}
