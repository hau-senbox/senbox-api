package usecase

import (
	"errors"
	"gorm.io/gorm"
	"sen-global-api/config"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
)

type GetToDoListByQRCodeUseCase struct {
	*repository.ToDoRepository
	r      *sheet.Reader
	dbConn *gorm.DB
}

func NewGetToDoListByQRCodeUseCase(cfg config.AppConfig, db *gorm.DB, r *sheet.Reader) *GetToDoListByQRCodeUseCase {
	return &GetToDoListByQRCodeUseCase{
		ToDoRepository: &repository.ToDoRepository{},
		r:              r,
		dbConn:         db,
	}
}

func (c *GetToDoListByQRCodeUseCase) Execute(qrCode string) (entity.SToDo, error) {
	todo, err := c.ToDoRepository.FindById(qrCode, c.dbConn)
	if err != nil {
		return entity.SToDo{}, err
	}

	if todo == nil {
		return entity.SToDo{}, errors.New("todo not found")
	}

	if todo.Type == value.ToDoTypeAssign || todo.SpreadsheetID == "" {
		return c.ToDoRepository.GetToDoListByQRCode(qrCode, c.dbConn)
	}

	return c.getTasksByComposeTodo(*todo)
}

func (c *GetToDoListByQRCodeUseCase) getTasksByComposeTodo(todo entity.SToDo) (entity.SToDo, error) {
	res, err := c.r.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: todo.SpreadsheetID,
		ReadRange:     todo.SheetName + "!K11:K11",
	})

	if err != nil {
		monitor.SendMessageViaTelegram("Cannot retrieve todo name ", todo.ID, " err ", err.Error())
		return c.ToDoRepository.GetToDoListByQRCode(todo.ID, c.dbConn)
	}

	todoName := ""
	if len(res) > 0 {
		todoName = res[0][0].(string)
	}

	todo.Name = todoName
	_, err = c.ToDoRepository.Save(c.dbConn, &todo)

	return c.ToDoRepository.GetToDoListByQRCode(todo.ID, c.dbConn)
}
