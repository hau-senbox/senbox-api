package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateUserEntityUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *UpdateUserEntityUseCase) UpdateUserEntity(req request.UpdateUserEntityRequest) error {
	return receiver.UpdateUser(req)
}
