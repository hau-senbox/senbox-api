package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type GetFunctionClaimPermissionUseCase struct {
	*repository.FunctionClaimPermissionRepository
}

func (receiver *GetFunctionClaimPermissionUseCase) GetAllFunctionClaimPermission(roleClaimId int64) ([]entity.SFunctionClaimPermission, error) {
	return receiver.GetAllByFunctionClaim(roleClaimId)
}

func (receiver *GetFunctionClaimPermissionUseCase) GetFunctionClaimPermissionById(req request.GetFunctionClaimPermissionByIdRequest) (*entity.SFunctionClaimPermission, error) {
	return receiver.GetByID(req)
}

func (receiver *GetFunctionClaimPermissionUseCase) GetFunctionClaimPermissionByName(req request.GetFunctionClaimPermissionByNameRequest) (*entity.SFunctionClaimPermission, error) {
	return receiver.GetByName(req)
}
