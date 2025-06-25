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

	var deviceID *string
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			deviceID, err = receiver.RegisteringDeviceForUser(user, req)
			if err != nil {
				return nil, err
			}
		}
	}

	if deviceID == nil {
		deviceID = &req.DeviceUUID
	}
	device, err := receiver.FindDeviceByID(*deviceID)
	if err != nil {
		return nil, err
	}

	// defer receiver.saveNewDeviceToSyncSheet(user, device, req)

	return &device.ID, err
}

func (receiver *RegisterDeviceUseCase) Reserve(deviceID string, appVersion string) error {

	return nil
	// device, _ := receiver.DeviceRepository.FindDeviceByID(deviceID)
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

	// spreadsheetID := match[1]

	// rowNo := 0
	// uuids, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
	// 	SpreadsheetID: spreadsheetID,
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
	// 		if uuid[0].(string) == deviceID {
	// 			rowNo = rowNumber + 12
	// 			break
	// 		}
	// 	}
	// }

	// deviceData := make([][]interface{}, 0)
	// deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02")}) //Created At
	// deviceData = append(deviceData, []interface{}{deviceID})                        //Device ID
	// deviceData = append(deviceData, []interface{}{appVersion})

	// if rowNo == 0 {
	// 	log.Error(fmt.Sprintf("failed to find placeholder row in sync devices sheet https://docs.google.com/spreadsheets/d/%s", spreadsheetID))
	// 	_, err := receiver.Writer.WriteRanges(sheet.WriteRangeParams{
	// 		Range:     "Devices!K" + strconv.Itoa(len(uuids)+12),
	// 		Rows:      deviceData,
	// 		Dimension: "COLUMNS",
	// 	}, spreadsheetID)

	// 	return err
	// } else {
	// 	return errors.New("this device is already existing on sync devices sheet")
	// }
}
