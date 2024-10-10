package usecase

import (
	"fmt"
	"regexp"
	"sen-global-api/internal/data/repository"
	"strings"
)

type UpdateOutputSubmissionSettingUseCase struct {
	SettingRepository *repository.SettingRepository
}

func (receiver *UpdateOutputSubmissionSettingUseCase) UpdateSubmissionSetting(spreadsheetUrl string, sheetName string) error {
	if !strings.Contains(strings.ToLower(spreadsheetUrl), "drive.google.com/drive/folders") {
		return fmt.Errorf("invalid google drive url")
	}
	re := regexp.MustCompile(`.*/([^?]+)`)
	match := re.FindStringSubmatch(spreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}

	folderId := match[1]
	return receiver.SettingRepository.UpdateSubmissionSetting(folderId, sheetName)
}
