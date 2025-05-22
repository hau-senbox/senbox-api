package usecase

import (
	"sen-global-api/internal/data/repository"
)

type ApproveOrgFormApplicationUseCase struct {
	*repository.OrganizationRepository
}

func (receiver *ApproveOrgFormApplicationUseCase) ApproveOrgFromApplication(applicationID int64) error {
	return receiver.OrganizationRepository.ApproveOrgFormApplication(applicationID)
}
