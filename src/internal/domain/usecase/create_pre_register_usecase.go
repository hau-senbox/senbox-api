package usecase

import (
	"sen-global-api/internal/data/repository"
)

type CreatePreRegisterUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *CreatePreRegisterUseCase) CreatePreRegister(email string) error {
	return receiver.UserEntityRepository.CreatePreRegisterUser(email)
}
