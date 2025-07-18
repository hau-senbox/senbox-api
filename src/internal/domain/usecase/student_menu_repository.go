package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type StudentMenuUseCase struct {
	Repo *repository.StudentMenuRepository
}

func NewStudentMenuUseCase(repo *repository.StudentMenuRepository) *StudentMenuUseCase {
	return &StudentMenuUseCase{Repo: repo}
}

func (uc *StudentMenuUseCase) Create(menu *entity.StudentMenu) error {
	return uc.Repo.Create(menu)
}

func (uc *StudentMenuUseCase) BulkCreate(menus []entity.StudentMenu) error {
	return uc.Repo.BulkCreate(menus)
}

func (uc *StudentMenuUseCase) DeleteByStudentID(studentID string) error {
	return uc.Repo.DeleteByStudentID(studentID)
}

func (uc *StudentMenuUseCase) GetByStudentID(studentID string) ([]entity.StudentMenu, error) {
	return uc.Repo.GetByStudentID(studentID)
}
