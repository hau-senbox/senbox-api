package request

import "sen-global-api/internal/domain/value"

type GetSubmissionByConditionRequest struct {
	AtrValueString string `json:"atr_value_string" binding:"required"`
}

type AtrValueString struct {
<<<<<<< HEAD
	UserID   string
	Key      *string
	DB       *string
	TimeSort value.TimeSort
	Duration *value.TimeRange
	Quantity int
=======
	UserID       string
	Key          *string
	DB           *string
	TimeSort     value.TimeSort
	DateDuration *value.TimeRange
	Quantity     *string
>>>>>>> develop
}
