package usecase

import (
	"github.com/google/uuid"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
)

type CreateUserUseCase struct {
	UserRepository *repository.UserRepository
}

func (receiver CreateUserUseCase) CreateUser(params request.CreateUserRequest) (*entity.SUser, error) {
	user := entity.SUser{
		UserId:      uuid.New().String(),
		Username:    params.Username,
		Fullname:    params.Fullname,
		Phone:       "",
		Email:       params.Email,
		Address:     "",
		Job:         "",
		CountryCode: "",
		Password:    params.Password,
		Role:        value.User,
	}
	return receiver.UserRepository.SaveUser(&user)
}
