package usecase

import (
	"encoding/json"
	"sen-global-api/internal/data/repository"

	log "github.com/sirupsen/logrus"
)

type GetShowPicsQuestionDetailUseCase struct {
	*repository.QuestionRepository
}

func (receiver *GetShowPicsQuestionDetailUseCase) Execute(questionId string) (string, error) {
	question, err := receiver.FindById(questionId)
	if err != nil {
		return "", err
	}
	log.Debug(`Buttons`, question)

	type Att struct {
		PhotoUrl string `json:"value"`
	}

	var att = Att{}
	err = json.Unmarshal([]byte(question.Attributes), &att)
	if err != nil {
		return "", err
	}
	log.Debug(`PhotoUrl`, att.PhotoUrl)

	return att.PhotoUrl, nil
}
