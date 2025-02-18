package response

import "time"

type Messaging struct {
	Email          []string `json:"email"`
	Value3         []string `json:"value3"`
	ShowMessageBox bool     `json:"showMessageBox"`
}

type QuestionAttributes struct {
	Value                string              `json:"value"`
	Number               int                 `json:"number"`
	Steps                int                 `json:"steps"`
	Options              []map[string]string `json:"options"`
	ButtonConfigurations []map[string]string `json:"button_configurations"`
	Messaging            Messaging           `json:"messaging"`
}

type QuestionListData struct {
	QuestionId     string             `json:"question_id"`
	QuestionType   string             `json:"question_type"`
	Question       string             `json:"question"`
	Attributes     QuestionAttributes `json:"attributes"`
	Order          int                `json:"order"`
	AnswerRequired bool               `json:"answer_required"`
	Enabled        bool               `json:"enabled"`
}

type QuestionListResponseData struct {
	QuestionListData []QuestionListData `json:"questions" binding:"required"`
	DecryptPassword  string             `json:"decrypt_password"`
	FormName         string             `json:"form_name" binding:"required"`
}

type QuestionListResponse struct {
	Data QuestionListResponseData `json:"data"`
}

type AllQuestionsResponseData struct {
	QuestionId   string             `json:"question_id"`
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
