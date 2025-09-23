package request

type UploadLanguageSettingRequest struct {
	LanguageSettings []LanguageSettingRequest `json:"language_settings"`
	DeleteIDs        []uint                   `json:"delete_ids"`
}

type LanguageSettingRequest struct {
	ID                 *uint  `json:"id"`
	LangKey            string `json:"lang_key"`
	RegionKey          string `json:"region_key"`
	IsPublished        bool   `json:"is_published"`
	DeactivatedMessage string `json:"deactivated_message"`
}
