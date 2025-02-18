package usecase

import (
	"errors"
	"fmt"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/job"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type SyncDevicesUseCase struct {
	*repository.SettingRepository
	*repository.DeviceRepository
	*repository.DeviceFormDatasetRepository
	*repository.UserEntityRepository
	*repository.UserConfigRepository
	*sheet.Reader
	*sheet.Writer
	TimeMachine           *job.TimeMachine
	UserSpreadsheetReader *sheet.Reader
	UserSpreadsheetWriter *sheet.Writer
}

func (receiver *SyncDevicesUseCase) ImportDevices(req request.SyncDevicesRequest) error {
	monitor.SendMessageViaTelegram(fmt.Sprintf("[INFO][SYNC] Start import devices with interval %d", req.Interval))
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}
	if req.Interval == 0 {
		req.AutoImport = false
	}

	err := receiver.SettingRepository.SaveSyncDevicesSettings(req)
	if err != nil {
		return err
	}

	spreadsheetId := match[1]
	values, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     `Devices!K11:AW`,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	// log.Debug(values)
	devices := []entity.SDevice{}
	for rowNo, row := range values {
		if len(row) < 13 {
			continue
		}
		action := row[3].(string)
		if strings.ToLower(action) != "upload" {
			if strings.ToLower(action) == "delete" {
				receiver.eraseDeviceAtRow(rowNo+11, spreadsheetId)
			}
			continue
		}

		device, err := receiver.CreateDevice(row, rowNo+11)
		if err != nil {
			// log.Error(err)
			continue
		}
		devices = append(devices, device)
	}

	for _, device := range devices {
		err = receiver.DeviceRepository.SaveOrUpdateDevices([]entity.SDevice{device})
		if err != nil {
			log.Error(err)
			continue
		}
		if device.ID == "" {
			continue
		}

		defer receiver.fetchDatasets(device)

		rowNo, err := receiver.findFirstRow(device.ID, values, 11)
		if err != nil {
			log.Error(err)
			continue
		}
		if rowNo == 0 {
			continue
		}
		deviceData := make([][]interface{}, 0)
		deviceData = append(deviceData, []interface{}{"UPLOADED"})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
		_, err = receiver.Writer.UpdateRange(sheet.WriteRangeParams{
			Range:     "Devices!O" + strconv.Itoa(rowNo),
			Rows:      deviceData,
			Dimension: "COLUMNS",
		}, spreadsheetId)
		if err != nil {
			log.Error(err)
			continue
		}
	}

	if !req.AutoImport {
		receiver.TimeMachine.ScheduleSyncDevices(0)
	} else {
		receiver.TimeMachine.ScheduleSyncDevices(req.Interval)
	}

	return nil
}

func (receiver *SyncDevicesUseCase) SyncDevices(req request.SyncDevicesRequest) error {
	monitor.SendMessageViaTelegram(fmt.Sprintf("[INFO][SYNC] Start sync devices with interval %d", req.Interval))
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}
	if req.Interval == 0 {
		req.AutoImport = false
	}

	err := receiver.SettingRepository.SaveSyncDevicesSettings(req)
	if err != nil {
		return err
	}

	spreadsheetId := match[1]
	values, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     `Devices!K11:AW`,
	})
	if err != nil {
		log.Error(err)
		return err
	}

	log.Debug(values)
	devices := []entity.SDevice{}
	for rowNo, row := range values {
		if len(row) < 13 {
			continue
		}
		action := row[3].(string)
		if strings.ToLower(action) != "upload" {
			if strings.ToLower(action) == "delete" {
				receiver.eraseDeviceAtRow(rowNo+11, spreadsheetId)
			}
			continue
		}

		device, err := receiver.CreateDevice(row, rowNo+11)
		if err != nil {
			// log.Error(err)
			continue
		}
		devices = append(devices, device)
	}

	for _, device := range devices {
		err = receiver.DeviceRepository.SaveOrUpdateDevices([]entity.SDevice{device})
		if err != nil {
			log.Error(err)
			continue
		}
		if device.ID == "" {
			continue
		}

		defer receiver.fetchDatasets(device)

		rowNo, err := receiver.findFirstRow(device.ID, values, 11)
		if err != nil {
			log.Error(err)
			continue
		}
		if rowNo == 0 {
			continue
		}
		deviceData := make([][]interface{}, 0)
		deviceData = append(deviceData, []interface{}{"UPLOADED"})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{nil})
		deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
		_, err = receiver.Writer.UpdateRange(sheet.WriteRangeParams{
			Range:     "Devices!O" + strconv.Itoa(rowNo),
			Rows:      deviceData,
			Dimension: "COLUMNS",
		}, spreadsheetId)

		if err != nil {
			log.Error(err)
			continue
		}
	}

	return nil
}

func (receiver *SyncDevicesUseCase) CreateDevice(rawData []interface{}, rowNo int) (entity.SDevice, error) {
	if rawData[1].(string) == "" {
		return entity.SDevice{}, errors.New("device id is empty")
	}
	rawInput := rawData[10].(string)
	input := value.InfoInputTypeKeyboard
	switch strings.ToLower(rawInput) {
	case "scan":
		input = value.InfoInputTypeBarcode
	case "keyboard":
		input = value.InfoInputTypeKeyboard
	case "api":
		input = value.InfoInputTypeBackOffice
	default:
		input = value.InfoInputTypeBarcode
	}

	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(rawData[12].(string))

	outputSpreadsheetId := ""
	if len(match) >= 2 {
		// return entity.SDevice{}, fmt.Errorf("invalid spreadsheet url")
		outputSpreadsheetId = match[1]
	}

	status := value.DeviceModeSuspended
	if len(rawData) > 14 {
		_status, err := value.GetDeviceModeFromString(strings.ToLower(rawData[13].(string)))
		if err != nil {
			return entity.SDevice{}, err
		}
		status = _status
	}

	screenButtonType := value.ScreenButtonType_Scan
	screenButtonValue := ""
	if strings.ToLower(rawData[6].(string)) == "scan" {
		screenButtonType = value.ScreenButtonType_Scan
		screenButtonValue = rawData[11].(string)
	} else {
		re1 := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
		match1 := re1.FindStringSubmatch(rawData[11].(string))
		if len(match1) >= 2 {
			screenButtonType = value.ScreenButtonType_List
			screenButtonValue = rawData[11].(string)
		} else {
			screenButtonType = value.ScreenButtonType_Scan
			screenButtonValue = rawData[11].(string)
		}
	}

	message := ""
	if len(rawData) > 14 {
		message = rawData[14].(string)
	}

	teacherSpreadsheetId := ""
	if len(rawData) > 22 {
		teacherSpreadsheetUrl := rawData[22].(string)
		re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
		match := re.FindStringSubmatch(teacherSpreadsheetUrl)

		if len(match) > 1 {
			teacherSpreadsheetId = match[1]
		}
	}

	userDevices, err := receiver.UserEntityRepository.GetUserDeviceById(rawData[1].(string))

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.SDevice{}, err
		}
	}

	log.Info("userDevices: ", userDevices)
	for _, ud := range *userDevices {
		user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIdRequest{ID: ud.UserId.String()})
		if err != nil {
			if !errors.Is(err, gorm.ErrRecordNotFound) {
				return entity.SDevice{}, err
			}
		}

		var userConfig *entity.SUserConfig
		if user.UserConfig != nil {
			userConfig, err = receiver.UserConfigRepository.GetByID(uint(user.UserConfigID.Int64))
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					return entity.SDevice{}, err
				}

				userConfig = nil
			}

			// update config
			if userConfig != nil {
				err = receiver.UserConfigRepository.UpdateUserConfig(request.UpdateUserConfigRequest{
					ID:                   uint(userConfig.ID),
					TopButtonConfig:      screenButtonValue,
					StudentOutputSheetId: outputSpreadsheetId,
					TeacherOutputSheetId: teacherSpreadsheetId,
				})
				if err != nil {
					return entity.SDevice{}, err
				}

				usrConfigId := uint(userConfig.ID)
				receiver.UserEntityRepository.UpdateUser(request.UpdateUserEntityRequest{
					ID:         user.ID.String(),
					Username:   user.Username,
					UserConfig: &usrConfigId,
				})
			}
		}

		if userConfig == nil {
			// create config
			userConfigId, err := receiver.UserConfigRepository.CreateUserConfig(request.CreateUserConfigRequest{
				TopButtonConfig:      screenButtonValue,
				StudentOutputSheetId: outputSpreadsheetId,
				TeacherOutputSheetId: teacherSpreadsheetId,
			})
			if err != nil {
				return entity.SDevice{}, err
			}

			usrConfigId := uint(*userConfigId)
			receiver.UserEntityRepository.UpdateUser(request.UpdateUserEntityRequest{
				ID:         user.ID.String(),
				Username:   user.Username,
				UserConfig: &usrConfigId,
			})
		}
	}

	device := entity.SDevice{
		ID:                rawData[1].(string),
		DeviceName:        rawData[5].(string),
		ScreenButtonType:  screenButtonType,
		Status:            status,
		InputMode:         input,
		DeactivateMessage: message,
		ButtonUrl:         rawData[10].(string),
		RowNo:             rowNo,
	}
	return device, nil
}

func (receiver *SyncDevicesUseCase) findFirstRow(id string, values [][]interface{}, startRow int) (int, error) {
	rowNo := 0
	for rowindex, row := range values {
		if len(row) > 1 {
			if row[1].(string) == id {
				return rowindex + startRow, nil
			}
		}
	}
	return rowNo, errors.New("Cannot determine row number for device id: " + id)
}

func (receiver *SyncDevicesUseCase) eraseDeviceAtRow(rowNo int, spreadsheetID string) {
	resetData := [][]interface{}{}
	resetData = append(resetData, []interface{}{""})
	resetData = append(resetData, []interface{}{""})
	resetData = append(resetData, []interface{}{""})
	resetData = append(resetData, []interface{}{nil})
	resetData = append(resetData, []interface{}{""})
	resetData = append(resetData, []interface{}{nil})
	resetData = append(resetData, []interface{}{nil})
	resetData = append(resetData, []interface{}{""})
	resetData = append(resetData, []interface{}{""})
	resetData = append(resetData, []interface{}{""})
	resetData = append(resetData, []interface{}{nil})
	resetData = append(resetData, []interface{}{nil})
	resetData = append(resetData, []interface{}{""})
	resetData = append(resetData, []interface{}{nil})
	resetData = append(resetData, []interface{}{nil})
	resetData = append(resetData, []interface{}{""})
	resetData = append(resetData, []interface{}{nil})
	resetData = append(resetData, []interface{}{""})

	_, err := receiver.Writer.UpdateRange(sheet.WriteRangeParams{
		Range:     "Devices!K" + strconv.Itoa(rowNo),
		Dimension: "COLUMNS",
		Rows:      resetData,
	}, spreadsheetID)

	if err != nil {
		monitor.SendMessageViaTelegram(fmt.Sprintf("[DELETE DEVICE]Cannot erase device at row %d\n[Error]: %s", rowNo, err.Error()))
	}
}

func (receiver *SyncDevicesUseCase) fetchDatasets(device entity.SDevice) {

	userDevices, err := receiver.UserEntityRepository.GetUserDeviceById(device.ID)
	if err != nil {
		return
	}

	datasets := []entity.SDeviceFormDataset{}
	for _, userDevice := range *userDevices {
		user, err := receiver.UserEntityRepository.GetByID(request.GetUserEntityByIdRequest{ID: userDevice.UserId.String()})
		if err != nil {
			return
		}

		userConfig, err := receiver.UserConfigRepository.GetByID(uint(user.UserConfigID.Int64))
		if err != nil || userConfig == nil {
			return
		}

		re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
		match := re.FindStringSubmatch(userConfig.TopButtonConfig)

		if len(match) < 2 {
			return
		}

		spreadsheetId := match[1]

		setData, err := receiver.UserSpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
			SpreadsheetId: spreadsheetId,
			ReadRange:     `USER_FORM!K11:WW`,
		})

		if err != nil {
			log.Error(err)
			return
		}

		if len(setData) > 0 {
			for i := 1; i < len(setData)-1; i++ {
				setValues := make(map[string]string)
				if len(setData[i]) == 0 {
					continue
				}
				setValues["set"] = setData[i][0].(string)
				if setValues["set"] == "" {
					continue
				}
				for j := 1; j < len(setData[i]); j++ {
					if setData[i][j] != nil {
						setValues[strings.ToLower(setData[0][j].(string))] = setData[i][j].(string)
					}
				}

				dataset := entity.SDeviceFormDataset{
					ID:                       device.ID + "_" + setData[i][0].(string),
					DeviceId:                 device.ID,
					Set:                      setValues["set"],
					QuestionDate:             setValues[value.GetStringValue(value.QuestionDate)],
					QuestionTime:             setValues[value.GetStringValue(value.QuestionTime)],
					QuestionDateTime:         setValues[value.GetStringValue(value.QuestionDateTime)],
					QuestionDurationForward:  setValues[value.GetStringValue(value.QuestionDurationForward)],
					QuestionDurationBackward: setValues[value.GetStringValue(value.QuestionDurationBackward)],
					QuestionScale:            setValues[value.GetStringValue(value.QuestionScale)],
					QuestionQRCode:           setValues[value.GetStringValue(value.QuestionQRCode)],
					QuestionSelection:        setValues[value.GetStringValue(value.QuestionSelection)],
					QuestionText:             setValues[value.GetStringValue(value.QuestionText)],
					QuestionCount:            setValues[value.GetStringValue(value.QuestionCount)],
					QuestionNumber:           setValues[value.GetStringValue(value.QuestionNumber)],
					QuestionPhoto:            setValues[value.GetStringValue(value.QuestionPhoto)],
					QuestionMultipleChoice:   setValues[value.GetStringValue(value.QuestionMultipleChoice)],
					QuestionButtonCount:      setValues[value.GetStringValue(value.QuestionButtonCount)],
					QuestionSingleChoice:     setValues[value.GetStringValue(value.QuestionSingleChoice)],
					QuestionButtonList:       setValues[value.GetStringValue(value.QuestionButtonList)],
					QuestionMessageBox:       setValues[value.GetStringValue(value.QuestionMessageBox)],
					QuestionShowPic:          setValues[value.GetStringValue(value.QuestionShowPic)],
					QuestionButton:           setValues[value.GetStringValue(value.QuestionButton)],
					QuestionPlayVideo:        setValues[value.GetStringValue(value.QuestionPlayVideo)],
					QuestionQRCodeFront:      setValues[value.GetStringValue(value.QuestionQRCodeFront)],
					QuestionChoiceToggle:     setValues[value.GetStringValue(value.QuestionChoiceToggle)],
					QuestionSignature:        setValues[value.GetStringValue(value.QuestionSignature)],
					QuestionWeb:              setValues[value.GetStringValue(value.QuestionWeb)],
					QuestionDraggableList:    setValues[value.GetStringValue(value.QuestionDraggableList)],
					QuestionSendMessage:      setValues[value.GetStringValue(value.QuestionSendMessage)],
				}
				datasets = append(datasets, dataset)
			}
		}
	}

	receiver.DeviceFormDatasetRepository.Save(datasets)
}
