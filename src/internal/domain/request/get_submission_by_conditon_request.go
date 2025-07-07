package request

import "sen-global-api/internal/domain/value"

type GetSubmissionByConditionRequest struct {
	AtrValueString string `json:"atr_value_string" binding:"required"`
}

type AtrValueString struct {
	UserID      string
	QuestionKey *string
	QuestionDB  *string
	TimeSort    value.TimeSort
	Duration    *value.TimeRange
}
