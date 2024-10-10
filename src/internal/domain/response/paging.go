package response

type Pagination struct {
	Page      int   `json:"page"`
	Limit     int   `json:"limit"`
	TotalPage int   `json:"total_page"`
	Total     int64 `json:"total"`
}
