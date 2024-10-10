package usecase

import (
	"context"
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
	"os"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/sheet"
	"strings"
	"time"
)

type UpdateDeviceUseCase struct {
	*repository.DeviceRepository
	*repository.SettingRepository
	*sheet.Writer
}

func (receiver *UpdateDeviceUseCase) UpdateDevice(deviceId string, req request.UpdateDeviceRequest) (*entity.SDevice, error) {
	device, err := receiver.DeviceRepository.GetDeviceById(deviceId)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		device.DeviceName = *req.Name
	}
	if req.ButtonUrl != nil {
		device.ButtonUrl = *req.ButtonUrl
	}
	if req.Message != nil {
		device.Message = *req.Message
	}
	if req.UserInfo != nil {
		device.PrimaryUserInfo = req.UserInfo.UserInfo1
		device.SecondaryUserInfo = req.UserInfo.UserInfo2
		device.TertiaryUserInfo = req.UserInfo.UserInfo3
		device.InputMode = value.InfoInputTypeBackOffice
	}

	if req.ScreenButton != nil {
		type ScreenButtons struct {
			ButtonType  value.ButtonType `json:"button_type"`
			ButtonTitle string           `json:"button_title"`
		}

		bttType := value.ScreenButtonType_Scan
		if strings.ToLower(req.ScreenButton.ButtonType) == "list" {
			bttType = value.ScreenButtonType_List
		}
		device.ScreenButtonType = bttType
		device.ScreenButtonValue = req.ScreenButton.ButtonValue
	}

	if req.UserInfo != nil {
		if req.UserInfo.UserInfo1 != "" {
			device.PrimaryUserInfo = req.UserInfo.UserInfo1
		}
		if req.UserInfo.UserInfo2 != "" {
			device.SecondaryUserInfo = req.UserInfo.UserInfo2
		}

		preInitDevice, err := receiver.DeviceRepository.FindByUserInfo(req.UserInfo.UserInfo1, req.UserInfo.UserInfo2)
		if err != nil {
			log.Infof("No existing device found for %s and %s", req.UserInfo.UserInfo1, req.UserInfo.UserInfo2)
		}
		if preInitDevice != nil {
			if preInitDevice.DeviceId != device.DeviceId {
				//In this case the sheet belong to the user info 1 and 2 is already used by another device
				//So we don't need to create a new sheet, but reassign the sheet to the requested device.
				device.PrimaryUserInfo = preInitDevice.PrimaryUserInfo
				device.SecondaryUserInfo = preInitDevice.SecondaryUserInfo
				device.SpreadsheetId = preInitDevice.SpreadsheetId
				dv, err := receiver.DeviceRepository.UpdateDevice(device)
				if err != nil {
					return nil, err
				}
				return dv, nil
			} else {
				dv, err := receiver.DeviceRepository.UpdateDevice(device)
				if err != nil {
					return nil, err
				}
				return dv, nil
			}
		} else {
			//In this case, the requested device is now re-init with the new user info 1 and 2
			//So we need to create a new sheet for the device
			outputSettingsData, err := receiver.SettingRepository.GetOutputSettings()
			if err != nil {
				return nil, errors.New("failed to get output settings")
			}
			var outputSettings OutputSetting
			if outputSettingsData != nil {
				err = json.Unmarshal([]byte(outputSettingsData.Settings), &outputSettings)
				if err != nil {
					return nil, err
				}
			}

			summarySettingData, err := receiver.SettingRepository.GetSummarySettings()
			if err != nil {
				return nil, errors.New("failed to get summary settings")
			}
			var summarySetting SummarySetting
			if summarySettingData != nil {
				err = json.Unmarshal([]byte(summarySettingData.Settings), &summarySetting)
				if err != nil {
					return nil, err
				}
			}

			srv, err := drive.NewService(context.Background(), option.WithCredentialsFile("./credentials/google_service_account.json"))
			if err != nil {
				log.Debug("Unable to access Drive API:", err)
			}

			//In this case, no device has user info 1 and 2 that the same with th request.
			//So we need to create a new device
			templateFilePath := "./config/output_template.xlsx" // File you want to upload on your PC
			baseMimeType := "text/xlsx"                         // mimeType of file you want to upload

			file, err := os.Open(templateFilePath)
			if err != nil {
				log.Error("Error: %v", err)
				return nil, err
			}
			defer file.Close()
			f := &drive.File{
				Name:     req.UserInfo.UserInfo1 + "_" + req.UserInfo.UserInfo2 + ".xlsx",
				Parents:  []string{outputSettings.FolderId},
				MimeType: "application/vnd.google-apps.spreadsheet",
			}
			res, err := srv.Files.Create(f).Media(file, googleapi.ContentType(baseMimeType)).Do()
			if err != nil {
				log.Error("Error: %v", err)
				return nil, err
			}

			log.Debug(res.DriveId)
			var spreadsheetId = res.Id

			answerRow := make([][]interface{}, 0)
			answerRow = append(answerRow, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
			answerRow = append(answerRow, []interface{}{req.UserInfo.UserInfo1})
			answerRow = append(answerRow, []interface{}{req.UserInfo.UserInfo2})
			answerRow = append(answerRow, []interface{}{req.UserInfo.UserInfo3})
			answerRow = append(answerRow, []interface{}{"https://docs.google.com/spreadsheets/d/" + spreadsheetId})

			params := sheet.WriteRangeParams{
				Range:     "Submissions!A2",
				Dimension: "COLUMNS",
				Rows:      answerRow,
			}
			writtenRanges, err := receiver.Writer.WriteRanges(params, summarySetting.SpreadsheetId)
			if err != nil {
				return nil, err
			}
			log.Debug(writtenRanges)

			//Init Account Sheet
			infoRows := make([][]interface{}, 0)
			infoRows = append(infoRows, []interface{}{device.DeviceId})
			infoRows = append(infoRows, []interface{}{req.UserInfo.UserInfo1})
			infoRows = append(infoRows, []interface{}{req.UserInfo.UserInfo2})
			infoRows = append(infoRows, []interface{}{req.UserInfo.UserInfo3})
			accountSheetParams := sheet.WriteRangeParams{
				Range:     "Account!C1",
				Dimension: "ROWS",
				Rows:      infoRows,
			}
			writtenAccountRanges, err := receiver.Writer.WriteRanges(accountSheetParams, spreadsheetId)
			if err != nil {
				return nil, err
			}
			log.Debug(writtenAccountRanges)

			device.PrimaryUserInfo = req.UserInfo.UserInfo1
			device.SecondaryUserInfo = req.UserInfo.UserInfo2
			device.SpreadsheetId = spreadsheetId
			dv, err := receiver.DeviceRepository.UpdateDevice(device)
			if err != nil {
				return nil, err
			}
			return dv, nil
		}
	}

	dv, err := receiver.DeviceRepository.UpdateDevice(device)
	if err != nil {
		return nil, err
	}
	return dv, nil
}
