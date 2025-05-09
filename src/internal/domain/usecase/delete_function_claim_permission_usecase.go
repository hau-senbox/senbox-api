package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type DeleteFunctionClaimPermissionUseCase struct {
	*repository.FunctionClaimPermissionRepository
}

func (receiver *DeleteFunctionClaimPermissionUseCase) DeleteFunctionClaimPermission(req request.DeleteFunctionClaimPermissionRequest) error {
	return receiver.FunctionClaimPermissionRepository.DeleteFunctionClaimPermission(req)
}
