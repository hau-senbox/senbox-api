package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type CreateUserEntityUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *CreateUserEntityUseCase) CreateUserEntity(req request.CreateUserEntityRequest) error {
	return receiver.CreateUser(req)
}
