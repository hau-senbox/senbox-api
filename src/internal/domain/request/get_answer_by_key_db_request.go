package request

type GetAnswerByKeyAndDB struct {
	AtrValueString string `json:"atr_value_string" binding:"required"`
}
