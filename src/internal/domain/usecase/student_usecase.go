package usecase

import (
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/response"
)

type StudentApplicationUseCase struct {
	StudentAppRepo *repository.StudentApplicationRepository
}

func NewStudentApplicationUseCase(repo *repository.StudentApplicationRepository) *StudentApplicationUseCase {
	return &StudentApplicationUseCase{
		StudentAppRepo: repo,
	}
}

// Get all students
func (uc *StudentApplicationUseCase) GetAllStudents() ([]response.StudentResponse, error) {
	apps, err := uc.StudentAppRepo.GetAll()
	if err != nil {
		return nil, err
	}

	res := make([]response.StudentResponse, 0, len(apps))
	for _, a := range apps {
		res = append(res, response.StudentResponse{
			StudentID:   fmt.Sprintf("%d", a.ID),
			StudentName: a.StudentName,
		})
	}
	return res, nil
}
