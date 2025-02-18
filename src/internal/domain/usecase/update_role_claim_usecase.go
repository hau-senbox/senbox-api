package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateRoleClaimUseCase struct {
	*repository.RoleClaimRepository
}

func (receiver *UpdateRoleClaimUseCase) UpdateRoleClaim(req request.UpdateRoleClaimRequest) error {
	return receiver.RoleClaimRepository.UpdateRoleClaim(req)
}
