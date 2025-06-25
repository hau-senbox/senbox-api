package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetQuestionByIDUseCase struct {
	repository.QuestionRepository
}

func (receiver *GetQuestionByIDUseCase) GetQuestionByID(id string) (*entity.SQuestion, error) {
	return receiver.FindByID(id)
}
