package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateUserOrgInfoUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *UpdateUserOrgInfoUseCase) UpdateUserOrgInfo(req request.UpdateUserOrgInfoRequest) error {
	return receiver.UserEntityRepository.UpdateUserOrgInfo(req)
}
