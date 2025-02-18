package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type DeleteRoleUseCase struct {
	*repository.RoleRepository
}

func (receiver *DeleteRoleUseCase) DeleteRole(req request.DeleteRoleRequest) error {
	return receiver.RoleRepository.DeleteRole(req)
}
