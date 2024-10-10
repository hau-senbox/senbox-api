package parameters

import "sen-global-api/internal/domain/value"

type RawQuestion struct {
	QuestionId        string                  `json:"question_id"`
	Question          string                  `json:"question"`
	Type              string                  `json:"type"`
	Attributes        string                  `json:"attributes"`
	AdditionalOptions string                  `json:"additional_options"`
	Status            string                  `json:"status"`
	AnswerRequired    string                  `json:"answer_required"`
	RowNumber         int                     `json:"row_number"`
	EnableOnMobile    value.QuestionForMobile `json:"enable_on_mobile"`
	QuestionUniqueId  *string `json:"question_unique_id"`
}

type SaveFormParams struct {
	Note              string
	Name              string
	SpreadsheetUrl    string
	SpreadsheetId     string
	Password          string
	RawQuestions      []RawQuestion
	SubmissionType    value.SubmissionType
	SubmissionSheetId string
	SheetName         string
	OutputSheetName   string
	SyncStrategy      value.FormSyncStrategy
}
