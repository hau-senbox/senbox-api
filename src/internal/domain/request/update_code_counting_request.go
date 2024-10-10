package request

type UpdateCodeCountingRequest struct {
	ID      uint `json:"id" binding:"required"`
	ResetTo int  `json:"reset_to" binding:"required"`
}
