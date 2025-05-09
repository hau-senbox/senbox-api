package usecase

import (
	"fmt"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
)

type UpdateEmailHistorySettingUseCase struct {
	*repository.SettingRepository
}

func (receiver *UpdateEmailHistorySettingUseCase) Execute(req request.UpdateEmailHistorySettingsRequest) error {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}

	spreadsheetId := match[1]

	return receiver.UpdateEmaiHistorySetting(spreadsheetId)
}
