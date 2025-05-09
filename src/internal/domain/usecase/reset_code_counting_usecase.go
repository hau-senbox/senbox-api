package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"

	"gorm.io/gorm"
)

type ResetCodeCountingUseCase struct {
	*repository.CodeCountingRepository
	*gorm.DB
}

func NewResetCodeCountingUseCase(db *gorm.DB) *ResetCodeCountingUseCase {
	return &ResetCodeCountingUseCase{
		CodeCountingRepository: &repository.CodeCountingRepository{},
		DB:                     db,
	}
}

func (receiver *ResetCodeCountingUseCase) Execute(req request.ResetCodeCountingRequest) error {
	return receiver.ResetCodeCounting(req, receiver.DB)
}
