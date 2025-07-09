package response

import "time"

type QuestionAttributes struct {
	Value                string              `json:"value"`
	Number               int                 `json:"number"`
	Steps                int                 `json:"steps"`
	Options              []map[string]string `json:"options"`
	ButtonConfigurations []map[string]string `json:"button_configurations"`
}

type QuestionListData struct {
	QuestionID     string             `json:"question_id"`
	QuestionType   string             `json:"question_type"`
	Question       string             `json:"question"`
	Attributes     QuestionAttributes `json:"attributes"`
	Order          int                `json:"order"`
	AnswerRequired bool               `json:"answer_required"`
	AnswerRemember bool               `json:"answer_remember"`
	RememberValue  string             `json:"remember_value"`
	Enabled        bool               `json:"enabled"`
	QuestionKey    string             `json:"question_key"`
	QuestionDB     string             `json:"question_db"`
}

type QuestionListResponseData struct {
	QuestionListData []QuestionListData `json:"questions" binding:"required"`
	DecryptPassword  string             `json:"decrypt_password"`
	FormName         string             `json:"form_name" binding:"required"`
}

type QuestionListResponse struct {
	Data   QuestionListResponseData `json:"data"`
	FormId uint64                   `json:"form_id"`
}

type AllQuestionsResponseData struct {
	QuestionID   string             `json:"question_id"`
	QuestionType string             `json:"question_type"`
	Question     string             `json:"question"`
	Attributes   QuestionAttributes `json:"attributes"`
	Status       string             `json:"status"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type AllQuestionsResponse struct {
	Data []AllQuestionsResponseData `json:"data"`
}
