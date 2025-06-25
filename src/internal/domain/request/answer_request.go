package request

type QuestionAnswer struct {
	QuestionID string `json:"question_id" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

type AnswerQuestionsRequest struct {
	Answers []QuestionAnswer `json:"answers" binding:"required"`
}
