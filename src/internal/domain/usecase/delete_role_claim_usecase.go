package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type DeleteFunctionClaimUseCase struct {
	*repository.FunctionClaimRepository
}

func (receiver *DeleteFunctionClaimUseCase) DeleteFunctionClaim(req request.DeleteFunctionClaimRequest) error {
	return receiver.FunctionClaimRepository.DeleteFunctionClaim(req)
}
