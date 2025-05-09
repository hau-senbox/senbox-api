package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type DeleteUserAuthorizeUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *DeleteUserAuthorizeUseCase) DeleteUserAuthorize(req request.DeleteUserAuthorizeRequest) error {
	return receiver.UserEntityRepository.DeleteUserAuthorize(req)
}
