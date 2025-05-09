package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type CreateFunctionClaimUseCase struct {
	*repository.FunctionClaimRepository
}

func (receiver *CreateFunctionClaimUseCase) CreateFunctionClaim(req request.CreateFunctionClaimRequest) error {
	return receiver.FunctionClaimRepository.CreateFunctionClaim(req)
}

func (receiver *CreateFunctionClaimUseCase) CreateFunctionClaims(req request.CreateFunctionClaimsRequest) error {
	for _, functionClaim := range req.FunctionClaims {
		err := receiver.FunctionClaimRepository.CreateFunctionClaim(functionClaim)
		if err != nil {
			return err
		}
	}
	return nil
}
