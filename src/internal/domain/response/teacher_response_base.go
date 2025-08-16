package response

type TeacherResponseBase struct {
	TeacherID      string                   `json:"id"`
	UserID         string                   `json:"user_id"`
	TeacherName    string                   `json:"name"`
	Avatar         string                   `json:"avatar,omitempty"`
	AvatarURL      string                   `json:"avatar_url,omitempty"`
	QrFormProfile  string                   `json:"qr_form,omitempty"`
	Menus          []ComponentResponse      `json:"components"`
	IsUserBlock    bool                     `json:"is_block"`
	LanguageConfig *LanguagesConfigResponse `json:"language_config"`
}
