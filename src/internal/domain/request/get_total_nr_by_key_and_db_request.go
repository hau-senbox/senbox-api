package request

type GetTotalNrByKeyAndDbRequest struct {
	AtrValueString string `json:"atr_value_string" binding:"required"`
}
