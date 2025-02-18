package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateRoleUseCase struct {
	*repository.RoleRepository
}

func (receiver *UpdateRoleUseCase) UpdateRole(req request.UpdateRoleRequest) error {
	return receiver.RoleRepository.UpdateRole(req)
}
