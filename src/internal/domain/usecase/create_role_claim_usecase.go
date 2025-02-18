package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type CreateRoleClaimUseCase struct {
	*repository.RoleClaimRepository
}

func (receiver *CreateRoleClaimUseCase) CreateRoleClaim(req request.CreateRoleClaimRequest) error {
	return receiver.RoleClaimRepository.CreateRoleClaim(req)
}

func (receiver *CreateRoleClaimUseCase) CreateRoleClaims(req request.CreateRoleClaimsRequest) error {
	for _, roleClaim := range req.RoleClaims {
		err := receiver.RoleClaimRepository.CreateRoleClaim(roleClaim)
		if err != nil {
			return err
		}
	}
	return nil
}
