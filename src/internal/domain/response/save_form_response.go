package response

import "time"

type SaveFormResponseData struct {
	Id          uint64    `json:"id"`
	Spreadsheet string    `json:"spreadsheet_url"`
	Password    string    `json:"password"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type SaveFormResponse struct {
	Data SaveFormResponseData `json:"data"`
}
