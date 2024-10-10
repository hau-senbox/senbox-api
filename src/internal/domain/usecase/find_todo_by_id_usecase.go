package usecase

import (
	"gorm.io/gorm"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type FindTodoByIdUseCase struct {
	TodoRepository repository.ToDoRepository
	DB             *gorm.DB
}

func NewFindTodoByIdUseCase(db *gorm.DB) *FindTodoByIdUseCase {
	return &FindTodoByIdUseCase{
		TodoRepository: repository.ToDoRepository{},
		DB:             db,
	}
}

func (c *FindTodoByIdUseCase) Execute(id string) (*entity.SToDo, error) {
	return c.TodoRepository.FindById(id, c.DB)
}
