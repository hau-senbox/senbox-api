package parameters

import (
	"sen-global-api/internal/domain/value"
)

type RawQuestion struct {
	QuestionID        string                  `json:"question_id"`
	Question          string                  `json:"question"`
	Type              string                  `json:"type"`
	Attributes        string                  `json:"attributes"`
	AdditionalOptions string                  `json:"additional_options"`
	Status            string                  `json:"status"`
	AnswerRequired    string                  `json:"answer_required"`
	AnswerRemember    string                  `json:"answer_remember"`
	RowNumber         int                     `json:"row_number"`
	EnableOnMobile    value.QuestionForMobile `json:"enable_on_mobile"`
	QuestionUniqueID  *string                 `json:"question_unique_id"`
	Key               string                  `json:"key"`
	DB                string                  `json:"db"`
}

type SaveFormParams struct {
	Note           string
	Name           string
	SpreadsheetUrl string
	SpreadsheetID  string
	Password       string
	RawQuestions   []RawQuestion
	SheetName      string
}
