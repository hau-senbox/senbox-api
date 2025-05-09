package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sen-global-api/config"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type UpdateToDoTasksUseCase struct {
	db                *gorm.DB
	repository        *repository.ToDoRepository
	settingRepository *repository.SettingRepository
	SpreadsheetWriter *sheet.Writer
	SpreadsheetReader *sheet.Reader
}

func NewUpdateToDoTasksUseCase(cfg config.AppConfig, db *gorm.DB, reader *sheet.Reader, writer *sheet.Writer) *UpdateToDoTasksUseCase {
	return &UpdateToDoTasksUseCase{
		db:         db,
		repository: &repository.ToDoRepository{},
		settingRepository: &repository.SettingRepository{
			DBConn: db,
		},
		SpreadsheetWriter: writer,
		SpreadsheetReader: reader,
	}
}

func (c *UpdateToDoTasksUseCase) UpdateTask(req request.UpdateToDoTasksRequest) (entity.SToDo, error) {
	importTodoSetting, err := c.settingRepository.GetSyncToDosSettings()
	if err != nil {
		return entity.SToDo{}, err
	}

	outputSettingsData, err := c.settingRepository.GetOutputSettings()
	if err != nil {
		return entity.SToDo{}, err
	}
	todo, err := c.repository.FindById(req.QRCode, c.db)
	if err != nil {
		return entity.SToDo{}, err
	}
	if todo == nil {
		return entity.SToDo{}, errors.New("ToDo not found")
	}
	if todo.Type != value.ToDoTypeCompose {
		return entity.SToDo{}, errors.New("ToDo type is not compose")
	}

	var tasks = make([]entity.Task, 0)

	for index, t := range req.Tasks {
		tasks = append(tasks, entity.Task{
			Index:     index,
			Name:      t.Name,
			DueDate:   t.DueDate,
			Value:     t.Value,
			Selection: t.Selection,
			Selected:  t.Selected,
		})
	}

	todo.Name = req.Name
	todo.Tasks = datatypes.JSONType[entity.STasks]{Data: entity.STasks{Tasks: tasks}}

	pwd, err := os.Getwd()
	if err != nil {
		monitor.SendMessageViaTelegram(fmt.Sprintf("Error getting current directory: %s", err))
		return entity.SToDo{}, err
	}
	if todo.SpreadsheetID == "" {
		var outputSettings OutputSetting
		if outputSettingsData != nil {
			err = json.Unmarshal([]byte(outputSettingsData.Settings), &outputSettings)
			if err != nil {
				return entity.SToDo{}, err
			}
		}
		srv, err := drive.NewService(context.Background(),
			option.WithCredentialsFile(pwd+"/credentials/google_service_account.json"),
		)
		if err != nil {
			return entity.SToDo{}, err
		}

		templateFilePath := pwd + "/config/todo_template.xlsx"    // File you want to upload on your PC
		baseMimeType := "application/vnd.google-apps.spreadsheet" // mimeType of file you want to upload

		file, err := os.Open(templateFilePath)
		if err != nil {
			log.Errorf("Error: %v", err)
			return entity.SToDo{}, err
		}
		log.Debug("File: ", file.Name())
		defer file.Close()
		f := &drive.File{
			Name:     todo.ID + ".xlsx",
			Parents:  []string{outputSettings.FolderId},
			MimeType: "application/vnd.google-apps.spreadsheet",
		}
		res, err := srv.Files.Create(f).Media(file, googleapi.ContentType(baseMimeType)).Do()
		if err != nil {
			log.Error("Error: ", err)
			monitor.SendMessageViaTelegram("Failed to create spreadsheet for ToDo")
			return entity.SToDo{}, err
		}
		todo.SpreadsheetID = res.Id

		//Write to todo uploader spread sheet
		var importSetting ImportSetting
		err = json.Unmarshal([]byte(importTodoSetting.Settings), &importSetting)
		if err != nil {
			return entity.SToDo{}, err
		}

		re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
		match := re.FindStringSubmatch(importSetting.SpreadSheetUrl)

		if len(match) < 2 {
			return entity.SToDo{}, fmt.Errorf("invalid spreadsheet url in import todo setting")
		}

		todoUploaderSpreadsheetId := match[1]

		//Find row where todo id belong to
		readColumnsK, err := c.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
			SpreadsheetId: todoUploaderSpreadsheetId,
			ReadRange:     "TODOs!K12:K1000",
		})
		if err != nil {
			return entity.SToDo{}, err
		}

		rowIndex := 12
		for index, r := range readColumnsK {
			if len(r) > 0 && r[0].(string) == todo.ID {
				rowIndex += index
				break
			}
		}
		if rowIndex != 12 {
			row := make([][]interface{}, 0)
			row = append(row, []interface{}{"https://docs.google.com/spreadsheets/d/" + todo.SpreadsheetID})
			updateUploaderParams := sheet.WriteRangeParams{
				Range:     "TODOs!L" + strconv.Itoa(rowIndex),
				Dimension: "COLUMNS",
				Rows:      row,
			}
			_, err = c.SpreadsheetWriter.UpdateRange(updateUploaderParams, todoUploaderSpreadsheetId)
			if err != nil {
				return entity.SToDo{}, err
			}
		} else {
			log.Error("TODO ", todo.ID, " does not exist from the TODO uploader")
		}
	} else {
		_, err := c.SpreadsheetWriter.ClearRange(sheet.ClearRangeParams{
			SpreadsheetId: todo.SpreadsheetID,
			Range:         todo.SheetName + "!I13:V500",
		})
		if err != nil {
			return entity.SToDo{}, err
		}
	}

	todoItems := make([][]interface{}, 0)
	for index, t := range tasks {
		date, err := time.Parse("2006-01-02 15:04:05", t.DueDate)
		if err != nil {
			todoItems = append(todoItems, []interface{}{
				index,
				t.Name,
				t.DueDate,
				t.Value,
				t.Selection,
				t.Selected,
			})
		} else {
			todoItems = append(todoItems, []interface{}{
				index,
				t.Name,
				date.Format("1/2/2006 15:04"),
				t.Value,
				t.Selection,
				t.Selected,
			})
		}
	}

	params := sheet.WriteRangeParams{
		Range:     todo.SheetName + "!K13",
		Dimension: "ROWS",
		Rows:      todoItems,
	}

	_, err = c.SpreadsheetWriter.WriteRanges(params, todo.SpreadsheetID)
	if err != nil {
		return entity.SToDo{}, err
	}

	updateTodoNameParams := sheet.WriteRangeParams{
		Range:     todo.SheetName + "!K11",
		Dimension: "ROWS",
		Rows:      [][]interface{}{{todo.Name}},
	}
	_, err = c.SpreadsheetWriter.UpdateRange(updateTodoNameParams, todo.SpreadsheetID)
	if err != nil {
		return entity.SToDo{}, err
	}

	monitor.SendMessageViaTelegram("Todo: ", todo.ID, " at ", todo.SpreadsheetID, " has been updated ")

	return c.repository.Save(c.db, todo)
}
