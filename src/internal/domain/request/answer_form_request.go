package request

import "time"

type Messaging struct {
	Email        []string `json:"email" binding:"required"`
	Value3       []string `json:"value3" binding:"required"`
	MessageBox   *string  `json:"messageBox"`
	QuestionType string   `json:"questionType" binding:"required"`
}

type Answer struct {
	QuestionId string     `json:"question_id" binding:"required"`
	Answer     string     `json:"answer" binding:"required"`
	Messaging  *Messaging `json:"messaging"`
}

type AnswerFormRequest struct {
	Answers  []Answer  `json:"answers" binding:"required"`
	OpenedAt time.Time `json:"opened_at"`
}
