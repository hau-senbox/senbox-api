package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetOrganizationUseCase struct {
	*repository.OrganizationRepository
}

func (receiver *GetOrganizationUseCase) GetOrganizationById(id int64) (*entity.SOrganization, error) {
	return receiver.GetByID(id)
}

func (receiver *GetOrganizationUseCase) GetAllOrganization(user *entity.SUserEntity) ([]*entity.SOrganization, error) {
	return receiver.GetAll(user)
}

func (receiver *GetOrganizationUseCase) GetAllUserByOrganization(organizationID uint) ([]*entity.SUserOrg, error) {
	return receiver.OrganizationRepository.GetAllUserByOrganization(organizationID)
}
