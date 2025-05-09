package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetOrganizationUseCase struct {
	*repository.OrganizationRepository
}

func (receiver *GetOrganizationUseCase) GetOrganizationById(id uint) (*entity.SOrganization, error) {
	return receiver.GetByID(id)
}

func (receiver *GetOrganizationUseCase) GetAllOrganization() ([]*entity.SOrganization, error) {
	return receiver.GetAll()
}

func (receiver *GetOrganizationUseCase) GetAllUserByOrganization(organizationID uint) ([]*entity.SUserEntity, error) {
	return receiver.OrganizationRepository.GetAllUserByOrganization(organizationID)
}
