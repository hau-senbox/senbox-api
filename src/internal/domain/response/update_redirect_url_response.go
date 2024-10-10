package response

type UpdateRedirectUrlResponse struct {
	Data GetRedirectUrlListResponseData `json:"data" binding:"required"`
}
