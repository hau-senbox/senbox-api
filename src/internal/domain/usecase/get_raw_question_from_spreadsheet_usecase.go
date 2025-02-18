package usecase

import (
	"sen-global-api/internal/domain/parameters"
	"sen-global-api/pkg/sheet"
	"strings"

	log "github.com/sirupsen/logrus"
)

type GetRawQuestionFromSpreadsheetUseCase struct {
	SpreadsheetId     string
	SpreadsheetReader *sheet.Reader
}

func (receiver *GetRawQuestionFromSpreadsheetUseCase) GetRawQuestions() ([]parameters.RawQuestion, error) {
	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: receiver.SpreadsheetId,
		ReadRange:     "Questions!A2:F",
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var result = make([]parameters.RawQuestion, 0)
	for index, row := range values {
		if len(row) == 6 && cap(row) == 6 {
			item := parameters.RawQuestion{
				QuestionId:        row[0].(string),
				Question:          row[2].(string),
				Type:              row[1].(string),
				Attributes:        strings.ReplaceAll(row[3].(string), "\n", ""),
				AdditionalOptions: row[4].(string),
				Status:            row[5].(string),
				RowNumber:         index + 1,
			}
			result = append(result, item)
		}
	}

	return result, err
}
