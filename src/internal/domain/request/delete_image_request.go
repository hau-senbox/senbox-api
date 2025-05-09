package request

type DeleteImageRequest struct {
	Key string `json:"key" binding:"required"`
}
