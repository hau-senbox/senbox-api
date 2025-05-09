package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetQuestionByIdUseCase struct {
	repository.QuestionRepository
}

func (receiver *GetQuestionByIdUseCase) GetQuestionById(id string) (*entity.SQuestion, error) {
	return receiver.FindById(id)
}
