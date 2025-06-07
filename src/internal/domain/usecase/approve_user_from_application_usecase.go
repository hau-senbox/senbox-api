package usecase

import (
	"sen-global-api/internal/data/repository"
)

type ApproveUserFormApplicationUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *ApproveUserFormApplicationUseCase) ApproveTeacherFromApplication(applicationID int64) error {
	return receiver.UserEntityRepository.ApproveTeacherFormApplication(applicationID)
}

func (receiver *ApproveUserFormApplicationUseCase) ApproveStaffFromApplication(applicationID int64) error {
	return receiver.UserEntityRepository.ApproveStaffFormApplication(applicationID)
}

func (receiver *ApproveUserFormApplicationUseCase) ApproveStudentFromApplication(applicationID int64) error {
	return receiver.UserEntityRepository.ApproveStudentFormApplication(applicationID)
}
