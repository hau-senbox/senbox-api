package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type GetFunctionClaimPermissionUseCase struct {
	*repository.FunctionClaimPermissionRepository
}

func (receiver *GetFunctionClaimPermissionUseCase) GetAllFunctionClaimPermission(roleClaimID int64) ([]entity.SFunctionClaimPermission, error) {
	return receiver.GetAllByFunctionClaim(roleClaimID)
}

func (receiver *GetFunctionClaimPermissionUseCase) GetFunctionClaimPermissionByID(req request.GetFunctionClaimPermissionByIDRequest) (*entity.SFunctionClaimPermission, error) {
	return receiver.GetByID(req)
}

func (receiver *GetFunctionClaimPermissionUseCase) GetFunctionClaimPermissionByName(req request.GetFunctionClaimPermissionByNameRequest) (*entity.SFunctionClaimPermission, error) {
	return receiver.GetByName(req)
}
