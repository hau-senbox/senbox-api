package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type CreateOrgFormApplicationUseCase struct {
	*repository.OrganizationRepository
	*repository.UserEntityRepository
}

func (receiver *CreateOrgFormApplicationUseCase) CreateOrgFormApplication(req request.CreateOrgFormApplicationRequest) error {
	_, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIDRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	return receiver.OrganizationRepository.CreateOrgFormApplication(req)
}
