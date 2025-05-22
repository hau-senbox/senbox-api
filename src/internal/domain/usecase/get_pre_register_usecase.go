package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetPreRegisterUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *GetPreRegisterUseCase) GetAllPreRegisterUser() ([]*entity.SPreRegister, error) {
	return receiver.UserEntityRepository.GetAllPreRegisterUser()
}

func (receiver *GetPreRegisterUseCase) GetPreRegisterUserByEmail(email string) (*entity.SPreRegister, error) {
	return receiver.UserEntityRepository.GetPreRegisterUserByEmail(email)
}
