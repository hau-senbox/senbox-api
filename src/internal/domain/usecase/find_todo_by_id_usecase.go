package usecase

import (
	"gorm.io/gorm"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type FindTodoByIDUseCase struct {
	TodoRepository repository.ToDoRepository
	DB             *gorm.DB
}

func NewFindTodoByIDUseCase(db *gorm.DB) *FindTodoByIDUseCase {
	return &FindTodoByIDUseCase{
		TodoRepository: repository.ToDoRepository{},
		DB:             db,
	}
}

func (c *FindTodoByIDUseCase) Execute(id string) (*entity.SToDo, error) {
	return c.TodoRepository.FindByID(id, c.DB)
}
