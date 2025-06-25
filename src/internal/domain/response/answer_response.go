package response

type AnsweredItem struct {
	QuestionID string `json:"question_id"`
	Question   string `json:"question"`
	Answer     string `json:"answer"`
	AnsweredAt string `json:"answered_at"`
}

type SucceedAnswerResponse struct {
	Data []AnsweredItem `json:"data"`
}
