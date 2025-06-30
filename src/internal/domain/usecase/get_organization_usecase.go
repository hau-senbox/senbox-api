package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetOrganizationUseCase struct {
	*repository.OrganizationRepository
}

func (receiver *GetOrganizationUseCase) GetOrganizationByID(id string) (*entity.SOrganization, error) {
	return receiver.OrganizationRepository.GetByID(id)
}

func (receiver *GetOrganizationUseCase) GetByName(name string) (*entity.SOrganization, error) {
	return receiver.OrganizationRepository.GetByName(name)
}

func (receiver *GetOrganizationUseCase) GetAllOrganization(user *entity.SUserEntity) ([]*entity.SOrganization, error) {
	return receiver.OrganizationRepository.GetAll(user)
}

func (receiver *GetOrganizationUseCase) GetAllUserByOrganization(organizationID string) ([]*entity.SUserOrg, error) {
	return receiver.OrganizationRepository.GetAllUserByOrganization(organizationID)
}
