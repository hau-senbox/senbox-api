package response

import "time"

type CreateFormSucceedData struct {
	FormId    uint64    `json:"form_id"`
	FormName  string    `json:"form_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateFormSucceedResponse struct {
	Data CreateFormSucceedData `json:"data"`
}

type GetFormListResponseData struct {
	Id          uint64    `json:"id"`
	Spreadsheet string    `json:"spreadsheet_url"`
	Password    string    `json:"password"`
	Note        string    `json:"note"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type GetFormListResponse struct {
	Data   []GetFormListResponseData `json:"data"`
	Paging Pagination                `json:"paging"`
}

type UpdateFormResponse struct {
	Data GetFormListResponseData `json:"data"`
}
