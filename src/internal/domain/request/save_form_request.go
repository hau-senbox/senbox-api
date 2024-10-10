package request

type SaveFormRequest struct {
	Note           string `json:"note" binding:"required"`
	SpreadsheetUrl string `json:"spreadsheet_url" binding:"required"`
	Password       string `json:"password"`
}
