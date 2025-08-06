package request

type SyncDataRequest struct {
	SheetUrl     string   `json:"sheet_url"`
	SheetName    string   `json:"sheet_name"`
	FormNotes    []string `json:"form_notes"` // giữ nguyên
	FormNotesStr string   `json:"form_qrs"`   // không bind từ JSON
	IsAuto       bool     `json:"is_auto"`
}
