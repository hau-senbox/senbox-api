package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type CreateUserFormApplicationUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *CreateUserFormApplicationUseCase) CreateTeacherFormApplication(req request.CreateTeacherFormApplicationRequest) error {
	_, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIdRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	return receiver.UserEntityRepository.CreateTeacherFormApplication(req)
}

func (receiver *CreateUserFormApplicationUseCase) CreateStaffFormApplication(req request.CreateStaffFormApplicationRequest) error {
	_, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIdRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	return receiver.UserEntityRepository.CreateStaffFormApplication(req)
}

func (receiver *CreateUserFormApplicationUseCase) CreateStudentFormApplication(req request.CreateStudentFormApplicationRequest) error {
	_, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIdRequest{ID: req.UserID})
	if err != nil {
		return err
	}

	return receiver.UserEntityRepository.CreateStudentFormApplication(req)
}
