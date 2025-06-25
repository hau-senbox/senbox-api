package request

import "time"

type Answer struct {
	QuestionID string `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
	Remember   bool   `json:"remember"`
}

type AnswerFormRequest struct {
	Answers  []Answer  `json:"answers" binding:"required"`
	OpenedAt time.Time `json:"opened_at"`
}
