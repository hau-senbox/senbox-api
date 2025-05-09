package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateUserAuthorizeUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *UpdateUserAuthorizeUseCase) UpdateUserAuthorize(req request.UpdateUserAuthorizeRequest) error {
	return receiver.UserEntityRepository.UpdateUserAuthorize(req)
}
