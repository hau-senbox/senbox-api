package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/queue"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

var serialQueue = queue.New()

type RegisterDeviceUseCase struct {
	*repository.UserRepository
	*repository.DeviceRepository
	*repository.SessionRepository
	*repository.SettingRepository
	*sheet.Writer
	*sheet.Reader
}

func (receiver *RegisterDeviceUseCase) RegisterDevice(user *entity.SUserEntity, req request.RegisterDeviceRequest) (*string, error) {
	deviceId, err := receiver.DeviceRepository.RegisteringDeviceForUser(user, req)
	if err != nil {
		return nil, err
	}

	device, err := receiver.DeviceRepository.FindDeviceById(*deviceId)

	defer receiver.saveNewDeviceToSyncSheet(user, device, req)

	return deviceId, err
}

func (receiver *RegisterDeviceUseCase) saveNewDeviceToSyncSheet(user *entity.SUserEntity, device *entity.SDevice, req request.RegisterDeviceRequest) {
	setting, err := receiver.SettingRepository.GetSyncDevicesSettings()
	if err != nil {
		log.Error("failed to get sync devices settings")
		return
	}

	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint8  `json:"interval"`
	}
	var importSetting ImportSetting
	err = json.Unmarshal([]byte(setting.Settings), &importSetting)
	if err != nil {
		log.Error("failed to unmarshal sync devices settings")
		return
	}
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(importSetting.SpreadSheetUrl)

	if len(match) < 2 {
		log.Error("failed to parse spreadsheet id from sync devices sheet")
		return
	}

	spreadsheetId := match[1]

	//Find first empty row at column L
	serialQueue <- func() {
		holderData, err := receiver.findPlaceholderRow(spreadsheetId, *device)
		if err == nil {
			receiver.insertDeviceToSyncSheet(user, holderData, device, req, spreadsheetId)
		} else {
			receiver.appendDeviceToSyncSheet(user, device, spreadsheetId)
		}
	}
}

type screenButtons struct {
	ButtonType  value.ButtonType `json:"button_type"`
	ButtonTitle string           `json:"button_title"`
}

type placeholderData struct {
	RowNo             int
	screenBtn         screenButtons
	InputModel        value.UserInfoInputType
	Status            value.DeviceStatus
	DeactivateMessage string
}

func (p placeholderData) getScreenButtonType() value.ScreenButtonType {
	switch p.screenBtn.ButtonType {
	case value.ButtonTypeScan:
		return value.ScreenButtonType_Scan
	case value.ButtonTypeList:
		return value.ScreenButtonType_List
	}

	return value.ScreenButtonType_Scan
}

func (p placeholderData) getDeviceMode() value.DeviceMode {
	switch p.Status {
	case value.DeviceStatus_ModeT:
		return value.DeviceModeT
	case value.DeviceStatus_ModeS:
		return value.DeviceModeS
	case value.DeviceStatus_ModeP:
		return value.DeviceModeP
	case value.DeviceStatus_ModeL:
		return value.DeviceModeL
	case value.DeviceStatus_Suspend:
		return value.DeviceModeSuspended
	case value.DeviceStatus_Deactive:
	}
	return value.DeviceModeT
}

func (receiver *RegisterDeviceUseCase) findPlaceholderRow(spreadsheetId string, device entity.SDevice) (placeholderData, error) {
	monitor.LogGoogleAPIRequestInitDevice()
	values, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Devices!K11:Z",
	})
	if err != nil {
		log.Error("failed to read from sync devices sheet to find empty L cell")
		return placeholderData{}, errors.New("no placeholder row found")
	}
	placeholder := placeholderData{
		RowNo: 0,
	}
	existingRow := 0
	for i, row := range values {
		var sBtn screenButtons = screenButtons{
			ButtonType:  value.ButtonTypeScan,
			ButtonTitle: "",
		}
		var inputModel value.UserInfoInputType = value.UserInfoInputTypeKeyboard
		var dstatus value.DeviceStatus = value.DeviceStatus_ModeT
		var message string = ""
		if len(row) > 1 {
			//Trim whiltespace
			var placeholderUUIDString string = row[1].(string)
			placeholderUUIDString = strings.ReplaceAll(placeholderUUIDString, " ", "")
			if placeholderUUIDString == "" {
				if len(row) > 6 {
					rawInput := row[6].(string)
					switch strings.ToLower(rawInput) {
					case "scan":
						inputModel = value.UserInfoInputTypeBarcode
					case "keyboard":
						inputModel = value.UserInfoInputTypeKeyboard
					case "api":
						inputModel = value.UserInfoInputTypeBackOffice
					default:
						inputModel = value.UserInfoInputTypeBarcode
					}
				}

				if len(row) > 11 {

					if strings.ToLower(row[6].(string)) == "scan" {
						sBtn = screenButtons{
							ButtonType:  value.ButtonTypeScan,
							ButtonTitle: row[11].(string),
						}
					} else {
						bttType := value.ButtonTypeScan
						re1 := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
						match1 := re1.FindStringSubmatch(row[11].(string))
						if len(match1) >= 2 {
							bttType = value.ButtonTypeList
							sBtn = screenButtons{
								ButtonType:  bttType,
								ButtonTitle: row[11].(string),
							}
						} else {
							sBtn = screenButtons{
								ButtonType:  value.ButtonTypeScan,
								ButtonTitle: row[11].(string),
							}
						}
					}
				}

				if len(row) > 14 {
					_status, err := value.GetDeviceStatusFromString(row[13].(string))
					if err == nil {
						dstatus = _status
					}
				}

				if len(row) > 14 {
					message = row[14].(string)
				}

				if placeholder.RowNo == 0 {
					placeholder = placeholderData{
						RowNo:             11 + i,
						screenBtn:         sBtn,
						InputModel:        inputModel,
						Status:            dstatus,
						DeactivateMessage: message,
					}
				}
			} else if placeholderUUIDString == device.ID {
				existingRow = 11 + i
			}
		}
	}

	if existingRow != 0 {
		screenButtonValue := value.ButtonTypeScan
		switch device.ScreenButtonType {
		case value.ScreenButtonType_Scan:
			screenButtonValue = value.ButtonTypeScan
		case value.ScreenButtonType_List:
			screenButtonValue = value.ButtonTypeList
		}
		var screenButtons = screenButtons{
			ButtonType:  screenButtonValue,
			ButtonTitle: "",
		}
		inputMode := value.UserInfoInputTypeBarcode
		switch device.InputMode {
		case value.InfoInputTypeKeyboard:
			inputMode = value.UserInfoInputTypeKeyboard
		case value.InfoInputTypeBarcode:
			inputMode = value.UserInfoInputTypeBarcode
		case value.InfoInputTypeBackOffice:
			inputMode = value.UserInfoInputTypeBackOffice
		}
		deviceStatus := value.DeviceStatus_ModeT
		switch device.Status {
		case value.DeviceModeSuspended:
			deviceStatus = value.DeviceStatus_Suspend
		case value.DeviceModeP:
			deviceStatus = value.DeviceStatus_ModeP
		case value.DeviceModeS:
			deviceStatus = value.DeviceStatus_ModeS
		case value.DeviceModeT:
			deviceStatus = value.DeviceStatus_ModeT
		case value.DeviceModeDeactivated:
			deviceStatus = value.DeviceStatus_Deactive
		}

		return placeholderData{
			RowNo:             existingRow,
			screenBtn:         screenButtons,
			InputModel:        inputMode,
			Status:            deviceStatus,
			DeactivateMessage: device.DeactivateMessage,
		}, nil
	} else if placeholder.RowNo != 0 {
		return placeholder, nil
	}

	return placeholderData{}, errors.New("no placeholder row found")
}

func (receiver *RegisterDeviceUseCase) appendDeviceToSyncSheet(user *entity.SUserEntity, device *entity.SDevice, spreadsheetId string) string {
	//Find first empty row
	var emptyRow int
	values, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Devices!K11:L1000",
	})
	if err != nil {
		log.Error("failed to read from sync devices sheet to find empty L cell")
		return ""
	}

	for i, row := range values {
		if len(row) == 0 {
			emptyRow = i + 11
			break
		}
	}
	if emptyRow == 0 {
		emptyRow = 11 + len(values)
	}

	var userGuardian *entity.SUserEntity
	if len(user.Guardians) > 0 {
		userGuardian = &user.Guardians[0]
	}

	guardianName := ""
	if userGuardian != nil {
		guardianName = userGuardian.Fullname
	}

	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{device.CreatedAt.Format("2006-01-02")})     //Created At
	deviceData = append(deviceData, []interface{}{device.ID})                                 //Device Id
	deviceData = append(deviceData, []interface{}{device.AppVersion})                         //Version
	deviceData = append(deviceData, []interface{}{nil})                                       //Command
	deviceData = append(deviceData, []interface{}{"UPLOADED"})                                //API Status
	deviceData = append(deviceData, []interface{}{nil})                                       //Device Name
	deviceData = append(deviceData, []interface{}{nil})                                       //Input Status
	deviceData = append(deviceData, []interface{}{user.Fullname})                             //User Info 1
	deviceData = append(deviceData, []interface{}{user.Company.CompanyName})                  //User Info 2
	deviceData = append(deviceData, []interface{}{guardianName})                              //User Info 3
	deviceData = append(deviceData, []interface{}{nil})                                       //Button Title., Screen Button setting
	deviceData = append(deviceData, []interface{}{nil})                                       //Button Url,App Sheet setting
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/"}) //App Sheet URL
	deviceData = append(deviceData, []interface{}{nil})                                       //Status
	deviceData = append(deviceData, []interface{}{nil})                                       //Message
	deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/"})

	monitor.LogGoogleAPIRequestInitDevice()
	response, err := receiver.Writer.UpdateRange(sheet.WriteRangeParams{
		Range:     "Devices!K" + strconv.Itoa(emptyRow),
		Rows:      deviceData,
		Dimension: "COLUMNS",
	}, spreadsheetId)
	if err != nil {
		log.Error("failed to write to sync devices sheet ", err.Error())
		return err.Error()
	}
	return response.UpdatedRange
}

func (receiver *RegisterDeviceUseCase) insertDeviceToSyncSheet(user *entity.SUserEntity, data placeholderData, device *entity.SDevice, req request.RegisterDeviceRequest, spreadsheetId string) {
	device.DeactivateMessage = data.DeactivateMessage
	device.Status = data.getDeviceMode()
	device.ButtonUrl = data.screenBtn.ButtonTitle
	device.ScreenButtonType = data.getScreenButtonType()
	device.InputMode = value.GetInfoInputTypeFromString(req.InputMode)
	device.RowNo = data.RowNo
	err := receiver.DeviceRepository.SaveDevices([]entity.SDevice{*device})
	if err != nil {
		log.Error("failed to save device")
		return
	}

	var userGuardian *entity.SUserEntity
	if len(user.Guardians) > 0 {
		userGuardian = &user.Guardians[0]
	}

	guardianName := ""
	if userGuardian != nil {
		guardianName = userGuardian.Fullname
	}

	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{device.CreatedAt.Format("2006-01-02")})     //Created At
	deviceData = append(deviceData, []interface{}{device.ID})                                 //Device Id
	deviceData = append(deviceData, []interface{}{device.AppVersion})                         //Version
	deviceData = append(deviceData, []interface{}{nil})                                       //Command
	deviceData = append(deviceData, []interface{}{"UPLOADED"})                                //API Status
	deviceData = append(deviceData, []interface{}{nil})                                       //Device Name
	deviceData = append(deviceData, []interface{}{nil})                                       //Input Status
	deviceData = append(deviceData, []interface{}{user.Fullname})                             //User Info 1
	deviceData = append(deviceData, []interface{}{user.Company.CompanyName})                  //User Info 2
	deviceData = append(deviceData, []interface{}{guardianName})                              //User Info 3
	deviceData = append(deviceData, []interface{}{nil})                                       //Button Title., Screen Button setting
	deviceData = append(deviceData, []interface{}{nil})                                       //Button Url,App Sheet setting
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/"}) //App Sheet URL
	deviceData = append(deviceData, []interface{}{nil})                                       //Status
	deviceData = append(deviceData, []interface{}{nil})                                       //Message
	deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/"})

	monitor.LogGoogleAPIRequestInitDevice()

	_, err = receiver.Writer.UpdateRange(sheet.WriteRangeParams{
		Range:     "Devices!K" + strconv.Itoa(data.RowNo),
		Rows:      deviceData,
		Dimension: "COLUMNS",
	}, spreadsheetId)

	if err != nil {
		log.Error("failed to write to sync devices sheet")
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][AUTHORIZE] Failed to insert into an placeholder row in sync devices sheet. Error: %s", err.Error()),
			fmt.Sprintf("Device Id: %s", device.ID),
		)
		return
	}

	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(data.screenBtn.ButtonTitle)

	if len(match) < 2 {
		log.Error("failed to get spreadsheet id to log accounts")
		return
	}

	accountSpreadsheetId := match[1]

	//Init Account Sheet
	infoRows := make([][]interface{}, 0)
	infoRows = append(infoRows, []interface{}{req.DeviceUUID})
	infoRows = append(infoRows, []interface{}{user.Fullname})
	infoRows = append(infoRows, []interface{}{user.Company.CompanyName})
	infoRows = append(infoRows, []interface{}{guardianName})
	infoRows = append(infoRows, []interface{}{req.AppVersion})
	accountSheetParams := sheet.WriteRangeParams{
		Range:     "Account!M11",
		Dimension: "ROWS",
		Rows:      infoRows,
	}
	monitor.LogGoogleAPIRequestInitDevice()
	writtenAccountRanges, err := receiver.Writer.UpdateRange(accountSheetParams, accountSpreadsheetId)
	if err != nil {
		log.Error("failed to write to account sheet")
		return
	} else {
		// monitor.SendMessageViaTelegram(
		// 	"[AUTHORIZE][insertDeviceToSyncSheet] Successfully updated account sheet",
		// 	fmt.Sprintf("Device ID: %s", device.ID),
		// 	fmt.Sprintf("Aargument: %v", infoRows),
		// 	fmt.Sprintf("Account sheet range: %v", writtenAccountRanges),
		// 	fmt.Sprintf("Account sheet ID: %s", accountSpreadsheetId),
		// )
	}
	log.Debug(writtenAccountRanges)
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
	// 	monitor.SendMessageViaTelegram(
	// 		"[ERROR][RESERVING] Cannot determine the row No of the device in sync devices sheet for reserve",
	// 		fmt.Sprintf("Device ID: %s is existing in the database", deviceId),
	// 		fmt.Sprintf("[Google sheet API error] %s", err.Error()),
	// 	)
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
