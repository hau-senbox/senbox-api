package usecase

import (
	"sen-global-api/internal/data/repository"
)

type BlockOrgFormApplicationUseCase struct {
	*repository.OrganizationRepository
}

func (receiver *BlockOrgFormApplicationUseCase) BlockOrgFormApplication(applicationID int64) error {
	return receiver.OrganizationRepository.BlockOrgFormApplication(applicationID)
}
