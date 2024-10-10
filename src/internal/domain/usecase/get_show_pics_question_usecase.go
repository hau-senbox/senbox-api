package usecase

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"sen-global-api/internal/data/repository"
)

type GetShowPicsQuestionDetailUseCase struct {
	*repository.QuestionRepository
}

func (receiver *GetShowPicsQuestionDetailUseCase) Execute(questionId string) (string, error) {
	question, err := receiver.QuestionRepository.FindById(questionId)
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
