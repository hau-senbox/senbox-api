package response

type GetRedirectUrlResponse struct {
	Data GetRedirectUrlListResponseData `json:"data" binding:"required"`
}
