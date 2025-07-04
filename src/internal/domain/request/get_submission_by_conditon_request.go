package request

import "sen-global-api/internal/domain/value"

type GetSubmissionByConditionRequest struct {
	AtrValueString string `json:"atr_value_string"`
}

type AtrValueString struct {
	QuestionKey string
	QuestionDB  string
	TimeSort    value.TimeSort
}
