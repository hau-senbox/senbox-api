package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sen-global-api/internal/domain/entity"
)

type ToDoRepository struct {
}

func (r *ToDoRepository) Save(conn *gorm.DB, list *entity.SToDo) (entity.SToDo, error) {
	conn.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "spreadsheet_id", "sheet_name", "tasks", "history_spreadsheet_id", "history_sheet_name", "updated_at", "start_row"}),
	}).Create(list)
	return *list, nil
}

func (r *ToDoRepository) GetToDoListByQRCode(code string, dbConn *gorm.DB) (entity.SToDo, error) {
	var todo entity.SToDo
	dbConn.Where("id = ?", code).First(&todo)

	if todo.ID == "" {
		return todo, gorm.ErrRecordNotFound
	}

	return todo, nil
}

func (r *ToDoRepository) FindById(id string, db *gorm.DB) (*entity.SToDo, error) {
	var todo entity.SToDo
	err := db.Where("id = ?", id).First(&todo).Error

	return &todo, err
}
