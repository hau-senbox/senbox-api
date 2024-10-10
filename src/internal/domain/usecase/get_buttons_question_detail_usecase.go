package usecase

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/sheet"
)

type GetButtonsQuestionDetailUseCase struct {
	*repository.QuestionRepository
	*sheet.Reader
}

func (receiver *GetButtonsQuestionDetailUseCase) Execute(questionId string) ([]response.GetScreenButtonsItem, error) {
	question, err := receiver.QuestionRepository.FindById(questionId)
	if err != nil {
		return nil, err
	}
	log.Debug(`Buttons`, question)

	type Att struct {
		SpreadsheetId string `json:"spreadsheet_id"`
	}

	var att = Att{}
	err = json.Unmarshal([]byte(question.Attributes), &att)
	if err != nil {
		return nil, err
	}
	log.Debug(`PhotoUrl`, att.SpreadsheetId)

	values, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: att.SpreadsheetId,
		ReadRange:     "Buttons!K12:M",
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var buttons []response.GetScreenButtonsItem
	for _, row := range values {
		if len(row) >= 2 {
			var title = row[0].(string)
			var value = row[1].(string)
			var backgroundColor string = ""
			if len(row) > 2 {
				backgroundColor = row[2].(string)
			}
			if backgroundColor != "" {
				buttons = append(buttons, response.GetScreenButtonsItem{
					ButtonTitle:     title,
					ButtonValue:     value,
					BackgroundColor: &backgroundColor,
				})
			} else {
				buttons = append(buttons, response.GetScreenButtonsItem{
					ButtonTitle:     title,
					ButtonValue:     value,
					BackgroundColor: nil,
				})
			}
		}
	}

	return buttons, nil
}
