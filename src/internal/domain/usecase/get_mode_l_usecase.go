package usecase

import (
	"errors"
	"regexp"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/pkg/sheet"

	log "github.com/sirupsen/logrus"
)

type GetModeLUseCase struct {
	Reader *sheet.Reader
	Writer *sheet.Writer
}

func (receiver *GetModeLUseCase) Execute(user entity.SUserEntity) (string, error) {
	if user.UserConfig == nil {
		log.Error("user config is nil")
		return "", errors.New("user config is nil")
	}

	topButtonConfig := user.UserConfig.TopButtonConfig
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(topButtonConfig)

	if len(match) < 2 {
		log.Error("failed to get spreadsheet id to log accounts")
		return "", errors.New("failed to get spreadsheet id of the account to find the mode of L")
	}

	accountSpreadsheetId := match[1]

	//Find the modeLValue of the cell M:19
	modeLValue, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: accountSpreadsheetId,
		ReadRange:     "Account!M19",
	})

	if err != nil {
		log.Error("failed to get the modeLValue of the cell M:19")
		return "", err
	}

	if len(modeLValue) < 1 {
		log.Error("failed to get the modeLValue of the cell M:19")
		return "", errors.New("failed to get the modeLValue of the cell M:19")
	}

	return modeLValue[0][0].(string), nil
}
