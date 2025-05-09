package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateFunctionClaimPermissionUseCase struct {
	*repository.FunctionClaimPermissionRepository
}

func (receiver *UpdateFunctionClaimPermissionUseCase) UpdateFunctionClaimPermission(req request.UpdateFunctionClaimPermissionRequest) error {
	return receiver.FunctionClaimPermissionRepository.UpdateFunctionClaimPermission(req)
}
