package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetAllQuestionsUseCase struct {
	*repository.QuestionRepository
}

func (receiver *GetAllQuestionsUseCase) GetAllQuestions() ([]entity.SQuestion, error) {
	return receiver.QuestionRepository.GetAllQuestions()
}
