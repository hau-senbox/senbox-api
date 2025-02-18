package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetUserConfigUseCase struct {
	*repository.UserConfigRepository
}

func (receiver *GetUserConfigUseCase) GetUserConfigById(id uint) (*entity.SUserConfig, error) {
	return receiver.UserConfigRepository.GetByID(id)
}
