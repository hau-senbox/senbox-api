package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateUserRoleUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *UpdateUserRoleUseCase) UpdateUserRole(req request.UpdateUserRoleRequest) error {
	return receiver.UserEntityRepository.UpdateUserRole(req)
}
