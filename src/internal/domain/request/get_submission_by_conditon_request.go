package request

import "sen-global-api/internal/domain/value"

type GetSubmissionByConditionRequest struct {
	ChildID        *string `json:"child_id"`
	AtrValueString string  `json:"atr_value_string" binding:"required"`
}

type AtrValueString struct {
	UserID       string
	Key          *string
	DB           *string
	TimeSort     value.TimeSort
	DateDuration *value.TimeRange
	Quantity     *string
}
