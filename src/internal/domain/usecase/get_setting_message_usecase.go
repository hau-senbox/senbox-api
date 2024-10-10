package usecase

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/pkg/sheet"
	"strconv"
)

type GetSettingMessageUseCase struct {
	*repository.DeviceRepository
	*sheet.Reader
}

func (receiver *GetSettingMessageUseCase) Execute(device entity.SDevice) (*response.GetSettingMessageResponse, error) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(device.ScreenButtonValue)

	if len(match) < 2 {
		return nil, errors.New("invalid spreadsheet url please contact BO")
	}

	spreadsheetId := match[1]

	values, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Time_Message!K11:O12",
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var fontSize *int = nil

	var items []response.GetSettingMessageItem

	for _, row := range values {
		if len(row) >= 2 && cap(row) >= 2 {
			items = append(items, response.GetSettingMessageItem{
				Description: row[0].(string),
				Message:     row[1].(string),
			})
		}

		if len(row) >= 5 && cap(row) >= 5 && fontSize == nil {
			rawFontSize := row[4].(string)
			if rawFontSize != "" {
				f, err := strconv.Atoi(rawFontSize)
				if f != 0 && err == nil {
					fontSize = &f
				}
			}
		}
	}

	return &response.GetSettingMessageResponse{
		Data: response.GetSettingMessageResponseData{
			Messages: items,
			FontSize: fontSize,
		},
	}, nil
}
