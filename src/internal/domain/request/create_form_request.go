package request

type CreateFormQuestionItem struct {
	QuestionId     string `json:"question_id"`
	Order          int    `json:"order"`
	AnswerRequired bool   `json:"answer_required"`
}

type CreateFormRequest struct {
	FormName  string                   `json:"form_name" binding:"required"`
	Questions []CreateFormQuestionItem `json:"questions" binding:"required"`
}
