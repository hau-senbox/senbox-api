package request

type GetChartNrRequest struct {
	AtrValueString string `json:"atr_value_string" binding:"required"`
}
