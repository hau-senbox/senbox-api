package usecase

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strings"
)

type UpdateApiDistributorUseCase struct {
	reader     *sheet.Reader
	writer     *sheet.Writer
	repository *repository.SettingRepository
}

func NewUpdateApiDistributorUseCase(db *gorm.DB, r *sheet.Reader, w *sheet.Writer) *UpdateApiDistributorUseCase {
	return &UpdateApiDistributorUseCase{
		reader:     r,
		writer:     w,
		repository: repository.NewSettingRepository(db),
	}
}

func (receiver *UpdateApiDistributorUseCase) Execute(url string) error {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(url)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}

	spreadsheetId := match[1]

	defer executeUploadAPIDistributor(receiver.repository, receiver.reader, receiver.writer)

	return receiver.repository.UpdateAPIDistributerSetting(spreadsheetId, url)
}

func executeUploadAPIDistributor(repo *repository.SettingRepository, r *sheet.Reader, w *sheet.Writer) {
	s, err := repo.GetAPIDistributerSetting()
	if err != nil {
		log.Error("Could not find API Distributor setting")
		monitor.SendMessageViaTelegram("Could not find API Distributor setting")
		return
	}

	var setting repository.APIDistributorSetting
	err = json.Unmarshal(s.Settings, &setting)
	if err != nil {
		log.Error("Invalid API setting record", err.Error())
		monitor.SendMessageViaTelegram("Invalid API setting record: ", err.Error())
		return
	}

	allSourceSheets, err := r.GetAllSheets(setting.SpreadSheetId)
	if err != nil {
		log.Error("Cannot retrieve all sheets from the api distributor spreadsheet")
		monitor.SendMessageViaTelegram("Cannot retrieve all sheets from the api distributor spreadsheet")
	}

	for _, singleSheet := range allSourceSheets {
		copyAPIDistributorAt(singleSheet, r, w, setting)
	}
}

func copyAPIDistributorAt(singleSheet sheet.SingleSheet, r *sheet.Reader, w *sheet.Writer, setting repository.APIDistributorSetting) {
	targets, err := r.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: setting.SpreadSheetId,
		ReadRange:     singleSheet.Title + "!M9:N",
	})

	if err != nil {
		log.Error("Cannot read sheet ", singleSheet.Title, " err ", err.Error())
		monitor.SendMessageViaTelegram("Cannot read sheet ", singleSheet.Title, " err ", err.Error())
		return
	}

	sourceSpreadsheetUrl := ""
	if len(targets) > 0 {
		if len(targets[0]) > 1 {
			sourceSpreadsheetUrl = targets[0][1].(string)
		}
	}

	if sourceSpreadsheetUrl == "" {
		log.Info("No source spreadsheet")
		return
	}

	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(sourceSpreadsheetUrl)

	if len(match) < 2 {
		log.Error(sourceSpreadsheetUrl, " is not a valid spreadsheet url")
		return
	}

	sourceSpreadsheetId := match[1]

	allSourceSheets, err := r.GetAllSheets(sourceSpreadsheetId)
	if err != nil {
		log.Error("Cannot fetch the sheets from source spreadsheet")
		return
	}

	var sourceSheet sheet.SingleSheet
	for _, tab := range allSourceSheets {
		if tab.Title == singleSheet.Title {
			sourceSheet = sheet.SingleSheet{
				ID:    tab.ID,
				Title: tab.Title,
			}
		}
	}

	log.Info("Source ", sourceSheet.Title)
	log.Info("Dis ", singleSheet.Title)

	if sourceSheet.Title == "" {
		log.Error("The sheet with name ", singleSheet.Title, " does not exist from source spreadsheet ", sourceSpreadsheetId)
		return
	}

	for _, row := range targets {
		if len(row) < 2 {
			continue
		}

		if strings.ToLower(row[0].(string)) != "upload" || row[1].(string) == "" {
			continue
		}

		re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
		match := re.FindStringSubmatch(row[1].(string))

		if len(match) < 2 {
			continue
		}

		targetSpreadsheetId := match[1]
		err = w.DeleteSheet(sheet.DeleteSheetParams{
			SpreadsheetID: targetSpreadsheetId,
			SheetTitle:    sourceSheet.Title,
		})

		if err != nil {
			log.Error("Failed on delete sheet ", singleSheet.Title+" from spreadsheet ", targetSpreadsheetId)
			monitor.SendMessageViaTelegram("Failed on delete sheet ", sourceSheet.Title+" from spreadsheet ", targetSpreadsheetId)
			continue
		}

		err = w.CopySingleSheet(sheet.CopySingleSheetParam{
			FromSpreadsheetId: sourceSpreadsheetId,
			SingleSheet:       sourceSheet,
			ToSpreadsheetId:   targetSpreadsheetId,
		})

		if err != nil {
			log.Error("Failed to copy sheet ", sourceSheet.Title, " to ", targetSpreadsheetId, " err ", err.Error())
			monitor.SendMessageViaTelegram("Failed to copy sheet ", singleSheet.Title, " to ", targetSpreadsheetId, " err ", err.Error())

			//TODO: update error in detail
			continue
		}
	}
}
