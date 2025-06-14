package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/pkg/sheet"

	"gorm.io/gorm"
)

type RegisterDeviceUseCase struct {
	*repository.DeviceRepository
	*repository.SessionRepository
	*repository.SettingRepository
	*sheet.Writer
	*sheet.Reader
}

func (receiver *RegisterDeviceUseCase) RegisterDevice(user *entity.SUserEntity, req request.RegisterDeviceRequest) (*string, error) {
	err := receiver.CheckUserDeviceExist(request.RegisteringDeviceForUser{
		UserID:   user.ID.String(),
		DeviceID: req.DeviceUUID,
	})

	var deviceId *string
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			deviceId, err = receiver.RegisteringDeviceForUser(user, req)
			if err != nil {
				return nil, err
			}
		}
	}

	if deviceId == nil {
		deviceId = &req.DeviceUUID
	}
	device, err := receiver.FindDeviceById(*deviceId)
	if err != nil {
		return nil, err
	}

	// defer receiver.saveNewDeviceToSyncSheet(user, device, req)

	return &device.ID, err
}

func (receiver *RegisterDeviceUseCase) Reserve(deviceId string, appVersion string) error {

	return nil
	// device, _ := receiver.DeviceRepository.FindDeviceById(deviceId)
	// if device != nil {
	// 	return errors.New("this device is already existing")
	// }

	// setting, err := receiver.SettingRepository.GetSyncDevicesSettings()
	// if err != nil {
	// 	log.Error("failed to get sync devices settings")
	// 	return err
	// }

	// type ImportSetting struct {
	// 	SpreadSheetUrl string `json:"spreadsheet_url"`
	// 	AutoImport     bool   `json:"auto"`
	// 	Interval       uint8  `json:"interval"`
	// }
	// var importSetting ImportSetting
	// err = json.Unmarshal([]byte(setting.Settings), &importSetting)
	// if err != nil {
	// 	log.Error("failed to unmarshal sync devices settings")
	// 	return err
	// }
	// re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	// match := re.FindStringSubmatch(importSetting.SpreadSheetUrl)

	// if len(match) < 2 {
	// 	log.Error("failed to parse spreadsheet id from sync devices sheet")
	// 	return err
	// }

	// spreadsheetId := match[1]

	// rowNo := 0
	// uuids, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
	// 	SpreadsheetId: spreadsheetId,
	// 	ReadRange:     "Devices!L12:L5000",
	// })

	// if err != nil {
	// 	log.Error("failed to find first row of sync devices sheet")
	// 	return err
	// }

	// firstEmptyRow := 0
	// for rowNumber, uuid := range uuids {
	// 	if len(uuid) == 0 && firstEmptyRow == 0 {
	// 		firstEmptyRow = rowNumber + 12
	// 		break
	// 	}
	// 	if len(uuid) != 0 {
	// 		if uuid[0].(string) == deviceId {
	// 			rowNo = rowNumber + 12
	// 			break
	// 		}
	// 	}
	// }

	// deviceData := make([][]interface{}, 0)
	// deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02")}) //Created At
	// deviceData = append(deviceData, []interface{}{deviceId})                        //Device Id
	// deviceData = append(deviceData, []interface{}{appVersion})

	// if rowNo == 0 {
	// 	log.Error(fmt.Sprintf("failed to find placeholder row in sync devices sheet https://docs.google.com/spreadsheets/d/%s", spreadsheetId))
	// 	_, err := receiver.Writer.WriteRanges(sheet.WriteRangeParams{
	// 		Range:     "Devices!K" + strconv.Itoa(len(uuids)+12),
	// 		Rows:      deviceData,
	// 		Dimension: "COLUMNS",
	// 	}, spreadsheetId)

	// 	return err
	// } else {
	// 	return errors.New("this device is already existing on sync devices sheet")
	// }
}
