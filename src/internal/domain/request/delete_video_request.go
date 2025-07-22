package request

type DeleteVideoRequest struct {
	Key string `json:"key" binding:"required"`
}
