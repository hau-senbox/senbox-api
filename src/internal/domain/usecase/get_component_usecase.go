package usecase

import (
	"sen-global-api/internal/data/repository"
)

type GetComponentUseCase struct {
	ComponentRepository *repository.ComponentRepository
}

func (receiver *GetComponentUseCase) GetAllComponentKey() (*[]string, error) {
	return receiver.ComponentRepository.GetAllComponentKey()
}
