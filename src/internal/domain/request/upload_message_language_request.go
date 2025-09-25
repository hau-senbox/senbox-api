package request

type UploadMessageLanguageRequest struct {
	TypeID     string `json:"type_id" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Key        string `json:"key" binding:"required"`
	Value      string `json:"message" binding:"required"`
	LanguageID uint   `json:"language_id" binding:"required"`
}

type UploadMessageLanguagesRequest struct {
	MessageLanguages []UploadMessageLanguageRequest `json:"message_languages" binding:"required"`
}
