package request

type UpdateOutputSubmissionSettingsRequest struct {
	FolderUrl string `json:"folder_url" binding:"required"`
	SheetName string `json:"sheet_name"`
}
