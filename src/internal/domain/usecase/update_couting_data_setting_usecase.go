package usecase

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

var DBConn *gorm.DB

func UpdateCodeCountingDataUseCase(req request.UpdateCodeCountingSettingRequest) error {
	var repo = repository.NewSettingRepository(DBConn)

	if req.SpreadsheetUrl == "" {
		return errors.New("code_counting_data_setting cannot be empty")
	}

	//Extract spreadsheet id from url
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		log.Error("failed to parse spreadsheet id from sync devices sheet")
		return errors.New("failed to parse spreadsheet id from sync devices sheet")
	}

	spreadsheetId := match[1]

	return repo.UpdateCodeCountingDataSetting(spreadsheetId, req.SpreadsheetUrl)
}
