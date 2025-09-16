package request

type UploadLanguageSettingRequest struct {
	LanguageSettings []LanguageSettingRequest `json:"language_settings"`
	DeleteIDs        []string                 `json:"delete_ids"`
}

type LanguageSettingRequest struct {
	ID        *uint  `json:"id"`
	LangKey   string `json:"lang_key" binding:"required"`
	RegionKey string `json:"region_key" binding:"required"`
}
