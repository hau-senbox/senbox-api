package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type CreateRoleUseCase struct {
	*repository.RoleRepository
}

func (receiver *CreateRoleUseCase) Create(req request.CreateRoleRequest) error {
	return receiver.CreateRole(req)
}
