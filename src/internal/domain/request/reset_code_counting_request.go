package request

type ResetCodeCountingRequest struct {
	Prefix  string `json:"prefix" binding:"required"`
	ResetTo int    `json:"reset_to" binding:"required"`
}
