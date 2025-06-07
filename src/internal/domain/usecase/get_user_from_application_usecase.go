package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type GetUserFormApplicationUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *GetUserFormApplicationUseCase) GetAllTeacherFromApplication() ([]*entity.STeacherFormApplication, error) {
	return receiver.UserEntityRepository.GetAllTeacherFormApplication()
}

func (receiver *GetUserFormApplicationUseCase) GetTeacherFromApplicationByID(applicationID int64) (*entity.STeacherFormApplication, error) {
	return receiver.UserEntityRepository.GetTeacherFormApplicationByID(applicationID)
}

func (receiver *GetUserFormApplicationUseCase) GetAllStaffFromApplication() ([]*entity.SStaffFormApplication, error) {
	return receiver.UserEntityRepository.GetAllStaffFormApplication()
}

func (receiver *GetUserFormApplicationUseCase) GetStaffFromApplicationByID(applicationID int64) (*entity.SStaffFormApplication, error) {
	return receiver.UserEntityRepository.GetStaffFormApplicationByID(applicationID)
}

func (receiver *GetUserFormApplicationUseCase) GetAllStudentFromApplication() ([]*entity.SStudentFormApplication, error) {
	return receiver.UserEntityRepository.GetAllStudentFormApplication()
}

func (receiver *GetUserFormApplicationUseCase) GetStudentFromApplicationByID(applicationID int64) (*entity.SStudentFormApplication, error) {
	return receiver.UserEntityRepository.GetStudentFormApplicationByID(applicationID)
}
