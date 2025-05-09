package usecase

import (
	"encoding/json"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"strings"
)

type GetUserQuestionsUseCase struct {
	repository.DeviceQuestionRepository
}

func (c GetUserQuestionsUseCase) GetQuestionsBelongToDevice(device *entity.SDevice) ([]response.QuestionListData, error) {
	userQuestions, err := c.DeviceQuestionRepository.GetQuestionsBelongToDevice(device)
	if err != nil {
		return nil, err
	}

	var result = make([]response.QuestionListData, len(userQuestions))
	for i, question := range userQuestions {
		var att response.QuestionAttributes
		_ = json.Unmarshal([]byte(question.Attributes), &att)
		result[i] = response.QuestionListData{
			QuestionId:   question.QuestionId,
			QuestionType: strings.ToUpper(question.QuestionType),
			Question:     question.Question,
			Attributes:   att,
			Order:        question.Order,
			Enabled:      question.EnableOnMobile == "enabled",
		}
	}

	return result, nil
}
