package request

type UploadMessageLanguageRequest struct {
	TypeID     string `json:"type_id" binding:"required"`
	Type       string `json:"type" binding:"required"`
	Key        string `json:"key" binding:"required"`
	Value      string `json:"message" binding:"required"`
	LanguageID uint   `json:"language_id" binding:"required"`
}

type UploadMessageLanguageListRequest struct {
	Messages []UploadMessageLanguageRequest `json:"messages" binding:"required,dive,required"`
}
