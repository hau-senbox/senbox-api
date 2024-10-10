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

type GetTimeTableUseCase struct {
	*repository.DeviceRepository
	*sheet.Reader
}

func (receiver *GetTimeTableUseCase) Execute(device entity.SDevice) (*response.GetTimeTableResponse, error) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(device.ScreenButtonValue)

	if len(match) < 2 {
		return nil, errors.New("invalid spreadsheet url please contact BO")
	}

	spreadsheetId := match[1]

	values, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Time_Message!K12:AG",
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var items []response.GetTimeTableItem
	var message string = ""
	var numberOfItemsPerTime int = 0

	for index, row := range values {

		if len(row) == 23 {
			re := regexp.MustCompile("[0-9]+")
			nums := re.FindAllString(row[21].(string), -1)
			if len(nums) > 0 {
				lastNum, err := strconv.Atoi(nums[len(nums)-1])
				if err == nil && numberOfItemsPerTime == 0 {
					numberOfItemsPerTime = lastNum
				}
			}
		}

		if index == 0 && len(row) >= 2 {
			message = row[1].(string)
		}
		if len(row) >= 2 && index > 3 {
			var color string = ""
			var message string = ""
			var notification string = ""
			var picture string = ""
			var link string = ""
			if len(row) >= 3 {
				color = row[2].(string)
			}
			if len(row) >= 4 {
				message = row[3].(string)
			}

			if len(row) >= 20 {
				notification = row[19].(string)
			}

			if len(row) >= 22 {
				picture = row[21].(string)
			}

			if len(row) >= 23 {
				link = row[22].(string)
			}

			items = append(items, response.GetTimeTableItem{
				StartAt:      row[0].(string),
				EndAt:        row[1].(string),
				Color:        color,
				Message:      message,
				Notification: notification,
				Picture:      picture,
				Link:         link,
			})
		}
	}

	return &response.GetTimeTableResponse{
		Data: response.GetTimeTableResponseData{
			GeneralMessage:       message,
			NumberOfItemsPerTime: numberOfItemsPerTime,
			Times:                items,
		},
	}, nil
}
