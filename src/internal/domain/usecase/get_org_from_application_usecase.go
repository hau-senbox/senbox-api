package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetOrgFormApplicationUseCase struct {
	*repository.OrganizationRepository
}

func (receiver *GetOrgFormApplicationUseCase) GetAllOrgFromApplication() ([]*entity.SOrgFormApplication, error) {
	return receiver.OrganizationRepository.GetAllOrgFormApplication()
}

func (receiver *GetOrgFormApplicationUseCase) GetOrgFromApplicationByID(applicationID int64) (*entity.SOrgFormApplication, error) {
	return receiver.OrganizationRepository.GetOrgFormApplicationByID(applicationID)
}
