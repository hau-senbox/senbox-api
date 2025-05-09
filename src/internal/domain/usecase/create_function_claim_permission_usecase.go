package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type CreateFunctionClaimPermissionUseCase struct {
	*repository.FunctionClaimPermissionRepository
}

func (receiver *CreateFunctionClaimPermissionUseCase) Create(req request.CreateFunctionClaimPermissionRequest) error {
	return receiver.CreateFunctionClaimPermission(req)
}
