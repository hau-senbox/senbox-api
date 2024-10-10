package usecase

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"regexp"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
)

type GetScreenButtonsByDeviceUseCase struct {
	Reader *sheet.Reader
}

func (receiver *GetScreenButtonsByDeviceUseCase) GetScreenButtons(device entity.SDevice) ([]response.GetScreenButtonsItem, error) {
	if device.ScreenButtonType == value.ScreenButtonType_Scan {
		return nil, errors.New("scan button not supported")
	}

	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(device.ScreenButtonValue)

	if len(match) < 2 {
		return nil, errors.New("invalid spreadsheet url please contact BO")
	}

	spreadsheetId := match[1]
	monitor.LogGoogleAPIRequestGETScreenButton()
	values, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Screen_Buttons!K12:M",
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	buttons := []response.GetScreenButtonsItem{}
	for _, row := range values {
		if len(row) >= 2 {
			var title = row[0].(string)
			var value = row[1].(string)
			var backgroundColor string = ""
			if len(row) > 2 {
				backgroundColor = row[2].(string)
			}
			if backgroundColor != "" {
				buttons = append(buttons, response.GetScreenButtonsItem{
					ButtonTitle:     title,
					ButtonValue:     value,
					BackgroundColor: &backgroundColor,
				})
			} else {
				buttons = append(buttons, response.GetScreenButtonsItem{
					ButtonTitle:     title,
					ButtonValue:     value,
					BackgroundColor: nil,
				})
			}
		}
	}

	return buttons, nil
}

func (receiver *GetScreenButtonsByDeviceUseCase) GetTopButtons(device entity.SDevice) ([]response.GetScreenButtonsItem, error) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(device.ScreenButtonValue)

	if len(match) < 2 {
		return nil, errors.New("invalid spreadsheet url please contact BO")
	}

	spreadsheetId := match[1]

	monitor.LogGoogleAPIRequestGETTopButton()
	values, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Top_Buttons!k12:M",
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	buttons := []response.GetScreenButtonsItem{}
	for _, row := range values {
		if len(row) >= 2 {
			var title = row[0].(string)
			var value = row[1].(string)
			var backgroundColor string = ""
			if len(row) > 2 {
				backgroundColor = row[2].(string)
			}
			if backgroundColor != "" {
				buttons = append(buttons, response.GetScreenButtonsItem{
					ButtonTitle:     title,
					ButtonValue:     value,
					BackgroundColor: &backgroundColor,
				})
			} else {
				buttons = append(buttons, response.GetScreenButtonsItem{
					ButtonTitle:     title,
					ButtonValue:     value,
					BackgroundColor: nil,
				})
			}
		}
	}

	return buttons, nil
}
