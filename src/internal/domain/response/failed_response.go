package response

type FailedResponse struct {
	Error Cause `json:"error"`
}

type Cause struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}
