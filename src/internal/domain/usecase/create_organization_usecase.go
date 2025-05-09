package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type CreateOrganizationUseCase struct {
	*repository.OrganizationRepository
}

func (receiver *CreateOrganizationUseCase) CreateOrganization(req request.CreateOrganizationRequest) error {
	return receiver.OrganizationRepository.CreateOrganization(req)
}
