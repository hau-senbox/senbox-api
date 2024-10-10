package request

type UpdateOutputTemplateRequest struct {
	SpreadsheetUrl string `json:"spreadsheet_url" binding:"required"`
}
