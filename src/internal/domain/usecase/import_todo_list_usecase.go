package usecase

import (
	"errors"
	"fmt"
	"regexp"
	"sen-global-api/config"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/job"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type ImportToDoListUseCase struct {
	cfg               config.AppConfig
	dbConn            *gorm.DB
	reader            *sheet.Reader
	writer            *sheet.Writer
	machine           *job.TimeMachine
	settingRepository *repository.SettingRepository
	todoRepository    *repository.ToDoRepository
}

func NewImportToDoListUseCase(cfg config.AppConfig, dbConn *gorm.DB, reader *sheet.Reader, writer *sheet.Writer, machine *job.TimeMachine) *ImportToDoListUseCase {
	return &ImportToDoListUseCase{
		cfg:               cfg,
		dbConn:            dbConn,
		reader:            reader,
		writer:            writer,
		machine:           machine,
		settingRepository: &repository.SettingRepository{DBConn: dbConn},
		todoRepository:    &repository.ToDoRepository{},
	}
}

func (receiver *ImportToDoListUseCase) ImportToDoList(req request.ImportFormRequest) error {
	monitor.SendMessageViaTelegram(fmt.Sprintf("[INFO][SYNC] Start sync ToDos with interval %d", req.Interval))
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}

	spreadsheetId := match[1]
	monitor.LogGoogleAPIRequestImportTodo()
	values, err := receiver.reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     `TODOs!` + receiver.cfg.Google.FirstColumn + strconv.Itoa(receiver.cfg.Google.FirstRow+2) + `:AA`,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	for rowNo, row := range values {
		err = receiver.saveToDoList(rowNo, row)
		if err != nil {
			log.Error(err)
		} else {
			log.Info("save todo list: ", rowNo)
			monitor.LogGoogleAPIRequestImportTodo()
			_, err = receiver.writer.UpdateRange(sheet.WriteRangeParams{
				Range:     "TODOs!P" + strconv.Itoa(rowNo+receiver.cfg.Google.FirstRow+2) + ":Q",
				Dimension: "ROWS",
				Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05")}},
			}, spreadsheetId)
			if err != nil {
				log.Debug("Row No: ", rowNo)
				log.Error(err)
			}
		}
	}

	err = receiver.settingRepository.UpdateSyncToDoSetting(req)
	if err != nil {
		return err
	}

	// if !req.AutoImport {
	// 	receiver.machine.ScheduleSyncToDos(0)
	// } else {
	// 	receiver.machine.ScheduleSyncToDos(req.Interval)
	// }

	return nil
}

func (receiver *ImportToDoListUseCase) saveToDoList(rowNo int, row []interface{}) error {
	if len(row) > 9 && strings.ToLower(row[4].(string)) != "upload" {
		if row[0].(string) != "" && !strings.Contains(strings.ToLower(row[0].(string)), "[todo-mobile]") && row[1].(string) != "" && strings.ToLower(row[4].(string)) == "upload" && row[9] != "" {
			re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
			match := re.FindStringSubmatch(row[1].(string))

			if len(match) < 2 {
				log.Info("invalid todo spreadsheet url")
				return errors.New("invalid todo spreadsheet url")
			}

			spreadsheetId := match[1]
			var tabName = "ToDos"
			if row[2].(string) != "" {
				tabName = row[2].(string)
			}

			var historySpreadsheetId string
			match = re.FindStringSubmatch(row[9].(string))
			if len(match) < 2 {
				log.Error("invalid history spreadsheet url")
				return errors.New("invalid history spreadsheet url")
			}
			historySpreadsheetId = match[1]

			return receiver.importToDo(spreadsheetId, tabName, row[0].(string), historySpreadsheetId)
		} else if row[0].(string) != "" && strings.Contains(strings.ToLower(row[0].(string)), "[todo-mobile]") && row[1].(string) == "" && strings.ToLower(row[4].(string)) == "upload" && row[9] != "" {
			re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
			match := re.FindStringSubmatch(row[9].(string))
			if len(match) < 2 {
				log.Error("invalid history spreadsheet url")
				return errors.New("invalid history spreadsheet url")
			}
			historySpreadsheetId := match[1]

			var tabName = "ToDos"
			if row[2].(string) != "" {
				tabName = row[2].(string)
			}

			return receiver.importToDoTypeCompose(row[0].(string), tabName, historySpreadsheetId)
		} else {
			log.Info("skip row: ", rowNo)
			return errors.New("skip row: " + strconv.Itoa(rowNo))
		}
	} else {
		log.Info("skip row: ", rowNo)
		return errors.New("skip row: " + strconv.Itoa(rowNo))
	}
}

func (receiver *ImportToDoListUseCase) importToDo(spreadsheetId string, tabName string, qrCode string, historySpreadsheetID string) error {
	monitor.LogGoogleAPIRequestImportTodo()
	values, err := receiver.reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     tabName + `!` + receiver.cfg.Google.FirstColumn + strconv.Itoa(11) + `:P`,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	var todoList = &entity.SToDo{
		ID:                   qrCode,
		Type:                 value.ToDoTypeAssign,
		SpreadsheetID:        spreadsheetId,
		SheetName:            tabName,
		HistorySpreadsheetID: historySpreadsheetID,
	}
	var tasks = make([]entity.Task, 0)
	startRow := 11
	todoName := ""
	for rowNum, row := range values {
		if len(row) == 1 && len(tasks) == 0 {
			todoName = row[0].(string)
		}
		if rowNum <= 1 {
			continue
		}
		if len(row) > 3 && !strings.Contains(strings.ToLower(row[1].(string)), "[todo]") {
			if row[0].(string) != "" && row[1].(string) != "" && row[2].(string) != "" {
				index, err := strconv.Atoi(row[0].(string))
				if err != nil {
					continue
				}
				dueInString := row[2].(string)
				if strings.ToLower(dueInString) == "urgent" || strings.ToLower(dueInString) == "no-date" || strings.ToLower(dueInString) == "form" {
					dueInString = strings.ToUpper(row[2].(string))
				} else {
					due, err := time.Parse("1/2/2006 15:04", row[2].(string))
					if err != nil {
						due, err = time.Parse("1/2/2006 15:04:05", row[2].(string))
						if err != nil {
							due, err = time.Parse("2 Jan 2006 - 15:04:05", row[2].(string))
							if err != nil {
								due, err = time.Parse("2 Jan 2006 - 15:04", row[2].(string))
								if err != nil {
									continue
								} else {
									dueInString = due.Format("2006-01-02 15:04:05")
								}
							} else {
								dueInString = due.Format("2006-01-02 15:04:05")
							}
						} else {
							dueInString = due.Format("2006-01-02 15:04:05")
						}
					} else {
						dueInString = due.Format("2006-01-02 15:04:05")
					}
				}

				selection := ""
				if len(row) > 4 {
					selection = row[4].(string)
				}
				selected := ""
				if len(row) > 5 {
					selected = row[5].(string)
				}
				tasks = append(tasks, entity.Task{
					Index:     index,
					Name:      row[1].(string),
					DueDate:   dueInString,
					Value:     row[3].(string),
					Selection: selection,
					Selected:  selected,
				})
				if startRow == 11 {
					startRow = receiver.cfg.Google.FirstRow + 1 + rowNum
				}
			}
		} else if len(row) > 1 {
			if strings.Contains(strings.ToLower(row[1].(string)), "[todo") {
				index, err := strconv.Atoi(row[0].(string))
				if err != nil {
					continue
				}
				tasks = append(tasks, entity.Task{
					Index:     index,
					Name:      row[1].(string),
					DueDate:   "",
					Value:     "",
					Selection: "",
					Selected:  "",
				})
				if startRow == 11 {
					startRow = receiver.cfg.Google.FirstRow + 1 + rowNum
				}
			}
		}
	}
	todoList.Name = todoName
	todoList.StartRow = startRow
	todoList.Tasks = datatypes.JSONType[entity.STasks]{Data: entity.STasks{Tasks: tasks}}

	_, _ = receiver.todoRepository.Save(receiver.dbConn, todoList)

	return nil
}

func (receiver *ImportToDoListUseCase) importToDoTypeCompose(qrCode, tabName, spreadsheetID string) error {
	monitor.LogGoogleAPIRequestImportTodo()
	var todoList = &entity.SToDo{
		ID:                   qrCode,
		Type:                 value.ToDoTypeCompose,
		SpreadsheetID:        "",
		SheetName:            tabName,
		HistorySpreadsheetID: spreadsheetID,
	}
	var tasks = make([]entity.Task, 0)

	todoList.Name = ""
	todoList.StartRow = 13
	todoList.Tasks = datatypes.JSONType[entity.STasks]{Data: entity.STasks{Tasks: tasks}}

	_, err := receiver.todoRepository.Save(receiver.dbConn, todoList)

	return err
}

func (receiver *ImportToDoListUseCase) ImportPartiallyToDos(spreadsheetURL string, sheetName string) error {
	monitor.SendMessageViaTelegram("[INFO][SYNC] Start import partially Todos")
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(spreadsheetURL)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}

	spreadsheetId := match[1]
	monitor.LogGoogleAPIRequestImportTodo()
	values, err := receiver.reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     sheetName + `!` + receiver.cfg.Google.FirstColumn + strconv.Itoa(receiver.cfg.Google.FirstRow+2) + `:AA`,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	for rowNo, row := range values {
		err = receiver.saveToDoList(rowNo, row)
		if err != nil {
			log.Error(err)
		} else {
			log.Info("save todo list: ", rowNo)
			monitor.LogGoogleAPIRequestImportTodo()
			_, err = receiver.writer.UpdateRange(sheet.WriteRangeParams{
				Range:     "TODOs!P" + strconv.Itoa(rowNo+receiver.cfg.Google.FirstRow+2) + ":Q",
				Dimension: "ROWS",
				Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05")}},
			}, spreadsheetId)
			if err != nil {
				log.Debug("Row No: ", rowNo)
				log.Error(err)
			}
		}
	}

	return nil
}
