package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type GetRoleUseCase struct {
	*repository.RoleRepository
}

func (receiver *GetRoleUseCase) GetAllRole() ([]entity.SRole, error) {
	return receiver.GetAll()
}

func (receiver *GetRoleUseCase) GetRoleByID(req request.GetRoleByIDRequest) (*entity.SRole, error) {
	return receiver.GetByID(req)
}

func (receiver *GetRoleUseCase) GetRoleByName(req request.GetRoleByNameRequest) (*entity.SRole, error) {
	return receiver.GetByName(req)
}
