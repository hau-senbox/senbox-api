package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type GetRoleClaimUseCase struct {
	*repository.RoleClaimRepository
}

func (receiver *GetRoleClaimUseCase) GetAllRoleClaimByRole(req request.GetAllRoleClaimByRoleRequest) ([]entity.SRoleClaim, error) {
	return receiver.GetAllByRole(req)
}

func (receiver *GetRoleClaimUseCase) GetAllRoleClaim() ([]entity.SRoleClaim, error) {
	return receiver.GetAll()
}

func (receiver *GetRoleClaimUseCase) GetRoleClaimById(req request.GetRoleClaimByIdRequest) (*entity.SRoleClaim, error) {
	return receiver.GetByID(req)
}

func (receiver *GetRoleClaimUseCase) GetRoleClaimByName(req request.GetRoleClaimByNameRequest) (*entity.SRoleClaim, error) {
	return receiver.GetByName(req)
}
