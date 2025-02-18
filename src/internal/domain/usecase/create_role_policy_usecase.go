package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type CreateRolePolicyUseCase struct {
	*repository.RolePolicyRepository
}

func (receiver *CreateRolePolicyUseCase) Create(req request.CreateRolePolicyRequest) error {
	return receiver.CreateRolePolicy(req)
}
