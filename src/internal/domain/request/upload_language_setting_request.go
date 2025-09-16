package request

type UploadLanguageSettingRequest struct {
	ID        *uint    `json:"id"`
	LangKey   string   `json:"lang_key" binding:"required"`
	RegionKey string   `json:"region_key" binding:"required"`
	DeleteIDs []string `json:"delete_ids"`
}
