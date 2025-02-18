package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type DeleteRolePolicyUseCase struct {
	*repository.RolePolicyRepository
}

func (receiver *DeleteRolePolicyUseCase) DeleteRolePolicy(req request.DeleteRolePolicyRequest) error {
	return receiver.RolePolicyRepository.DeleteRolePolicy(req)
}
