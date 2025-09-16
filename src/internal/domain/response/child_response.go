package response

import "sen-global-api/internal/domain/entity"

type ChildResponse struct {
	ChildID        string                   `json:"id"`
	ChildName      string                   `json:"name"`
	Avatar         string                   `json:"avatar"`
	AvatarURL      string                   `json:"avatar_url"`
	Parent         *entity.SUserEntity      `json:"parent,omitempty"`
	QrFormProfile  string                   `json:"qr_form"`
	Menus          []ComponentResponse      `json:"components"`
	LanguageConfig *LanguagesConfigResponse `json:"language_config"`
	Avatars        []Avatar                 `json:"avatars"`
	CreatedIndex   int                      `json:"created_index"`
}
