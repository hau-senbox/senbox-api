package request

type UpdateEmailHistorySettingsRequest struct {
	SpreadsheetUrl string `json:"spreadsheet_url" binding:"required"`
}
