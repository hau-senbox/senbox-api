package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"firebase.google.com/go/v4/messaging"
	"fmt"
	log "github.com/sirupsen/logrus"
	"gorm.io/datatypes"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"time"
)

func ModifyDevice(req request.ModifyDeviceRequest, deviceID string) error {
	log.Info("Start modify device with id: ", deviceID)

	repo := repository.NewDeviceRepository(DBConn)

	// Update device name
	if req.DeviceName != nil {
		err := repo.UpdateDeviceName(deviceID, *req.DeviceName)
		if err != nil {
			return err
		}
	}

	// Update Device Message
	if req.Message != nil {
		err := repo.UpdateDeviceMessage(deviceID, *req.Message)
		if err != nil {
			return err
		}
		defer func() {
			go announceDeviceUserMessageHasChanged(deviceID, *req.Message)
		}()
	}

	// Update Device Note
	if req.Note != nil {
		err := repo.UpdateDeviceNote(deviceID, *req.Note)
		if err != nil {
			return err
		}
		defer func() {
			go announceDeviceNoteHasChanged(deviceID, *req.Note)
		}()
	}

	// Update Device Status
	if req.Status != nil {
		_, err := value.GetDeviceStatusFromString(*req.Status)
		if err != nil {
			return err
		}
		err = repo.UpdateDeviceMode(deviceID, *req.Status)
		if err != nil {
			return err
		}

		defer func() {
			go announceDeviceStatusHasChanged(deviceID, *req.Status)
		}()
	}

	//App Setting Spreadsheet Url
	if req.AppSettingSpreadsheetUrl != nil {
		err := repo.UpdateAppSettingSpreadsheetUrl(deviceID, *req.AppSettingSpreadsheetUrl)
		if err != nil {
			return err
		}
	}

	//Output Spreadsheet Url
	if req.OutputSpreadsheetUrl != nil {
		err := repo.UpdateOutputSpreadsheetUrl(deviceID, *req.OutputSpreadsheetUrl)
		if err != nil {
			return err
		}
	}

	return nil
}

func announceDeviceStatusHasChanged(deviceID string, statusInString string) {
	repo := repository.NewMobileDeviceRepository()
	device, err := repo.FindByDeviceID(deviceID, DBConn)
	if err != nil {
		return
	}

	ctx := context.Background()
	msgApp, err := FirebaseApp.Messaging(ctx)
	if err != nil {
		monitor.SendMessageViaTelegram("[ERROR][INFORM DEVICE STATUS] Cannot initialize Messaging App ", err.Error())
		return
	}

	msg := &messaging.Message{
		Token: device.FCMToken,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Data: map[string]string{
				"status": statusInString,
				"type":   string(value.NotificationType_DeviceStatusChanged),
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
				},
				CustomData: map[string]interface{}{
					"status": statusInString,
					"type":   string(value.NotificationType_DeviceStatusChanged),
				},
			},
		},
	}

	_, err = msgApp.Send(ctx, msg)
	if err != nil {
		monitor.SendMessageViaTelegram("[ERROR][INFORM DEVICE STATUS] Cannot send message ", err.Error())
		return
	}

	monitor.SendMessageViaTelegram("Device status has changed to ", statusInString+" for device: "+deviceID)
}

func announceDeviceUserMessageHasChanged(deviceID, message string) {
	repo := repository.NewMobileDeviceRepository()
	device, err := repo.FindByDeviceID(deviceID, DBConn)
	if err != nil {
		return
	}

	ctx := context.Background()
	msgApp, err := FirebaseApp.Messaging(ctx)
	if err != nil {
		monitor.SendMessageViaTelegram("[ERROR][User-Message] Cannot initialize Messaging App ", err.Error())
		return
	}

	msg := &messaging.Message{
		Token: device.FCMToken,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Data: map[string]string{
				"message": message,
				"type":    string(value.NotificationType_UserMessageChanged),
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
				},
				CustomData: map[string]interface{}{
					"message": message,
					"type":    value.NotificationType_UserMessageChanged,
				},
			},
		},
	}

	_, err = msgApp.Send(ctx, msg)
	if err != nil {
		monitor.SendMessageViaTelegram("[ERROR][INFORM MESSAGE CHANGED] Cannot send notification", " error ", err.Error())
		return
	}

	monitor.SendMessageViaTelegram("User message has changed to ", message+" for device: "+deviceID)
}

func announceDeviceNoteHasChanged(deviceID, note string) {
	repo := repository.NewMobileDeviceRepository()
	device, err := repo.FindByDeviceID(deviceID, DBConn)
	if err != nil {
		return
	}

	ctx := context.Background()
	msgApp, err := FirebaseApp.Messaging(ctx)
	if err != nil {
		monitor.SendMessageViaTelegram("[ERROR][User-Message] Cannot initialize Messaging App ", err.Error())
		return
	}

	msg := &messaging.Message{
		Token: device.FCMToken,
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Data: map[string]string{
				"note": note,
				"type": string(value.NotificationType_NoteChanged),
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
				},
				CustomData: map[string]interface{}{
					"note": note,
					"type": value.NotificationType_NoteChanged,
				},
			},
		},
	}

	_, err = msgApp.Send(ctx, msg)
	if err != nil {
		monitor.SendMessageViaTelegram("[ERROR][INFORM MESSAGE CHANGED] Cannot send notification", " error ", err.Error())
		return
	}

	monitor.SendMessageViaTelegram("Note has changed to ", note+" for device: "+deviceID)
}

func SyncDevice(deviceID string) error {
	device, err := repository.FindDeviceByDeviceID(deviceID, DBConn)
	if err != nil {
		return err
	}

	if device.RowNo == 0 {
		return errors.New(fmt.Sprintf("Cannot find device's row for id %s in device uploader sheet, please contact web admin to insert it manually", device.DeviceId))
	}

	log.Info("SyncDevice: ", device.DeviceId, " RowNo: ", device.RowNo)
	//Find Device Sync Setting
	setting, err := repository.FindDeviceSyncSetting(DBConn)

	if err != nil {
		return err
	}

	spreadsheetID, err := extractSpreadsheetID(setting.Settings)
	if err != nil {
		return err
	}

	log.Info("SpreadsheetID: ", spreadsheetID)
	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{device.AppVersion})
	deviceData = append(deviceData, []interface{}{"UPLOAD"})
	deviceData = append(deviceData, []interface{}{"UPLOADED"})
	deviceData = append(deviceData, []interface{}{device.DeviceName})

	//Scan Or Press
	switch device.InputMode {
	case value.InfoInputTypeBarcode:
		deviceData = append(deviceData, []interface{}{"SCAN"})
		break
	case value.InfoInputTypeKeyboard:
		deviceData = append(deviceData, []interface{}{"PRESS"})
		break
	case value.InfoInputTypeBackOffice:
		deviceData = append(deviceData, []interface{}{""})
		break
	}

	//Values
	deviceData = append(deviceData, []interface{}{device.PrimaryUserInfo})
	deviceData = append(deviceData, []interface{}{device.SecondaryUserInfo})
	deviceData = append(deviceData, []interface{}{device.TertiaryUserInfo})

	deviceData = append(deviceData, []interface{}{device.ButtonUrl})
	deviceData = append(deviceData, []interface{}{device.ScreenButtonValue})
	deviceData = append(deviceData, []interface{}{fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", device.SpreadsheetId)})

	//Status
	switch device.Status {
	case value.DeviceModeL:
		deviceData = append(deviceData, []interface{}{"MODE L"})
		break
	case value.DeviceModeP:
		deviceData = append(deviceData, []interface{}{"MODE P"})
		break
	case value.DeviceModeS:
		deviceData = append(deviceData, []interface{}{"MODE S"})
		break
	case value.DeviceModeT:
		deviceData = append(deviceData, []interface{}{"MODE T"})
		break
	case value.DeviceModeSuspended:
		deviceData = append(deviceData, []interface{}{"SUSPEND"})
		break
	case value.DeviceModeDeactivated:
		deviceData = append(deviceData, []interface{}{"DEACTIVE"})
		break
	}

	deviceData = append(deviceData, []interface{}{device.Message})
	deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02 15:04:05")})

	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{fmt.Sprintf("https://docs.google.com/spreadsheets/d/%s", device.TeacherSpreadsheetId)})

	//Write to Spreadsheet
	writeParams := sheet.WriteRangeParams{
		Range:     "Devices!M" + strconv.Itoa(device.RowNo),
		Dimension: "COLUMNS",
		Rows:      deviceData,
	}

	_, err = AdminSpreadsheetClient.Writer.UpdateRange(writeParams, spreadsheetID)
	if err == nil {
		monitor.SendMessageViaTelegram("SyncDevice: ", device.DeviceId, " RowNo: ", strconv.Itoa(device.RowNo))
	}

	return err
}

func extractSpreadsheetID(settings datatypes.JSON) (string, error) {
	type ImportSetting struct {
		SpreadSheetUrl string `json:"spreadsheet_url"`
		AutoImport     bool   `json:"auto"`
		Interval       uint8  `json:"interval"`
	}
	var importSetting ImportSetting
	err := json.Unmarshal(settings, &importSetting)
	if err != nil {
		log.Error("failed to unmarshal sync devices settings")
		return "", err
	}

	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(importSetting.SpreadSheetUrl)

	if len(match) < 2 {
		log.Error("failed to parse spreadsheet id from sync devices sheet")
		return "", err
	}

	return match[1], nil
}

func GetDeviceList(request request.GetListDeviceRequest) ([]entity.SDevice, *response.Pagination, error) {
	repo := repository.NewDeviceRepository(DBConn)

	return repo.GetDeviceList(request)
}
