package request

type SetupLogoRefreshIntervalRequest struct {
	Interval uint64 `json:"interval"`
	Title    string `json:"title"`
}
