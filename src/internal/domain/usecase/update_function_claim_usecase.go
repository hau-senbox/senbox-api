package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateFunctionClaimUseCase struct {
	*repository.FunctionClaimRepository
}

func (receiver *UpdateFunctionClaimUseCase) UpdateFunctionClaim(req request.UpdateFunctionClaimRequest) error {
	return receiver.FunctionClaimRepository.UpdateFunctionClaim(req)
}
