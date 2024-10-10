package request

type UpdateCodeCountingSettingRequest struct {
	SpreadsheetUrl string `json:"spreadsheet_url" binding:"required"`
}
