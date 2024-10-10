package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"
)

type GetUsersUseCase struct {
	UserRepository *repository.UserRepository
}

func (receiver GetUsersUseCase) GetUsers(role value.Role) ([]*entity.SUser, error) {
	return receiver.UserRepository.GetUsers(role)
}
