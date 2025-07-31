package request

type SyncDataRequest struct {
	SheetUrl       string `json:"sheet_url"`
	SheetName      string `json:"sheet_name"`
	LastSubmitTime string `json:"last_submit_time"`
}
