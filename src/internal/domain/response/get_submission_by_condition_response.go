package response

import (
	"time"
)

type GetSubmissionByConditionResponse struct {
	SubmissionID uint64 `json:"id"`
	Key          string `json:"key"`
	DB           string `json:"db"`
	QuestionID   string `json:"question_id" binding:"required"`
	Question     QuestionListData
	Answer       string    `json:"answer" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
}

type SubmissionData struct {
	Items []GetSubmissionByConditionResponse `json:"items" binding:"required"`
}
