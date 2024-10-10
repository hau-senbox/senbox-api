package request

type ImportPartiallyTodoRequest struct {
	SpreadsheetURL string `json:"spreadsheet_url" binding:"required"`
	TabName        string `json:"tab_name" binding:"required"`
}
