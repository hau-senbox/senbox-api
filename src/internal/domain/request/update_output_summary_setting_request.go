package request

type UpdateOutputSummarySettingsRequest struct {
	SpreadsheetUrl string `json:"spreadsheet_url" binding:"required"`
}
