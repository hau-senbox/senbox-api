package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type DeleteRoleClaimUseCase struct {
	*repository.RoleClaimRepository
}

func (receiver *DeleteRoleClaimUseCase) DeleteRoleClaim(req request.DeleteRoleClaimRequest) error {
	return receiver.RoleClaimRepository.DeleteRoleClaim(req)
}
