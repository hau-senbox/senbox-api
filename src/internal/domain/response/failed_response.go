package response

type FailedResponse struct {
	Code    int    `json:"status_code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}
