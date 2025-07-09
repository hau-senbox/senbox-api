package request

type GetSubmission4MemmoriesRequest struct {
	AtrValueListString string `json:"atr_value_list_string" binding:"required"`
}
