package response

type SyncQueueResponse struct {
	ID        uint64 `json:"id"`
	SheetURL  string `json:"sheet_url"`
	SheetName string `json:"sheet_name"`
	FormQRs   string `json:"form_qrs"`
	IsAuto    *bool  `json:"is_auto"`
}
