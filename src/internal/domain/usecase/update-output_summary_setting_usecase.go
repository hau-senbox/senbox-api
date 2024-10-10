package usecase

import (
	"fmt"
	"regexp"
	"sen-global-api/internal/data/repository"
)

type UpdateOutputSummarySettingUseCase struct {
	SettingRepository *repository.SettingRepository
}

func (receiver *UpdateOutputSummarySettingUseCase) UpdateOutputSummarySetting(spreadsheetUrl string) error {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(spreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}

	spreadsheetId := match[1]

	return receiver.SettingRepository.UpdateOutputSummarySetting(spreadsheetId)
}
