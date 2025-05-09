package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type GetFunctionClaimUseCase struct {
	*repository.FunctionClaimRepository
}

func (receiver *GetFunctionClaimUseCase) GetAllFunctionClaim() ([]entity.SFunctionClaim, error) {
	return receiver.GetAll()
}

func (receiver *GetFunctionClaimUseCase) GetFunctionClaimById(req request.GetFunctionClaimByIdRequest) (*entity.SFunctionClaim, error) {
	return receiver.GetByID(req)
}

func (receiver *GetFunctionClaimUseCase) GetFunctionClaimByName(req request.GetFunctionClaimByNameRequest) (*entity.SFunctionClaim, error) {
	return receiver.GetByName(req)
}
