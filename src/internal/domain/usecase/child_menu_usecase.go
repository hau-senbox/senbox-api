package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type ChildMenuUseCase struct {
	Repo *repository.ChildMenuRepository
}

func NewChildMenuUseCase(repo *repository.ChildMenuRepository) *ChildMenuUseCase {
	return &ChildMenuUseCase{Repo: repo}
}

func (uc *ChildMenuUseCase) Create(menu *entity.ChildMenu) error {
	return uc.Repo.Create(menu)
}

func (uc *ChildMenuUseCase) BulkCreate(menus []entity.ChildMenu) error {
	return uc.Repo.BulkCreate(menus)
}

func (uc *ChildMenuUseCase) DeleteByChildID(childID string) error {
	return uc.Repo.DeleteByChildID(childID)
}

func (uc *ChildMenuUseCase) GetByChildID(childID string) ([]entity.ChildMenu, error) {
	return uc.Repo.GetByChildID(childID)
}
