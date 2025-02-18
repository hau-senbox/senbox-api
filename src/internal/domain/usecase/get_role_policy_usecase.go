package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type GetRolePolicyUseCase struct {
	*repository.RolePolicyRepository
}

func (receiver *GetRolePolicyUseCase) GetAllRolePolicy() ([]entity.SRolePolicy, error) {
	return receiver.GetAll()
}

func (receiver *GetRolePolicyUseCase) GetRolePolicyById(req request.GetRolePolicyByIdRequest) (*entity.SRolePolicy, error) {
	return receiver.GetByID(req)
}

func (receiver *GetRolePolicyUseCase) GetRolePolicyByName(req request.GetRolePolicyByNameRequest) (*entity.SRolePolicy, error) {
	return receiver.GetByName(req)
}
