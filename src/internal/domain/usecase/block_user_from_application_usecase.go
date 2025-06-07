package usecase

import (
	"sen-global-api/internal/data/repository"
)

type BlockUserFormApplicationUseCase struct {
	*repository.UserEntityRepository
}

func (receiver *BlockUserFormApplicationUseCase) BlockTeacherFormApplication(applicationID int64) error {
	return receiver.UserEntityRepository.BlockTeacherFormApplication(applicationID)
}

func (receiver *BlockUserFormApplicationUseCase) BlockStaffFormApplication(applicationID int64) error {
	return receiver.UserEntityRepository.BlockStaffFormApplication(applicationID)
}

func (receiver *BlockUserFormApplicationUseCase) BlockStudentFormApplication(applicationID int64) error {
	return receiver.UserEntityRepository.BlockStudentFormApplication(applicationID)
}
