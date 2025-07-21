package response

import "time"

type GetAnswerByKeyAndDbResponse struct {
	ID           string    `json:"id"`
	SubmissionID uint64    `json:"submission_id"`
	QuestionID   string    `json:"question_id"`
	Question     string    `json:"question"`
	UserID       string    `json:"user_id"`
	UserNickName string    `json:"user_nick_name"`
	Key          string    `json:"key"`
	DB           string    `json:"db"`
	Answer       string    `json:"answer" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
}
