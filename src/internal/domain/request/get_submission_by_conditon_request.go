package request

type GetSubmissionByConditionRequest struct {
	FormID      uint64 `json:"form_id" binding:"required"`
	UserID      string `json:"user_id"`
	QuestionKey string `json:"question_key"`
	QuestionDB  string `json:"question_db"`
}
