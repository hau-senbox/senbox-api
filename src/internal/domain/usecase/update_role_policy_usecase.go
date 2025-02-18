package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateRolePolicyUseCase struct {
	*repository.RolePolicyRepository
}

func (receiver *UpdateRolePolicyUseCase) UpdateRolePolicy(req request.UpdateRolePolicyRequest) error {
	return receiver.RolePolicyRepository.UpdateRolePolicy(req)
}
