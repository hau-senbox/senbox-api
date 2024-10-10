package request

type ImportRedirectUrlsRequest struct {
	SpreadsheetUrl string `json:"spreadsheet_url" binding:"required"`
	AutoImport     bool   `json:"auto"`
	Interval       uint64 `json:"interval_in_minutes"`
}
