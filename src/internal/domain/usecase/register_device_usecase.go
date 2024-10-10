package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/queue"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/googleapi"
	"google.golang.org/api/option"
)

type AuthorizeUsers struct {
	PrimaryUser   entity.SUser
	SecondaryUser entity.SUser
	TertiaryUser  entity.SUser
}

var serialQueue = queue.New()

type RegisterDeviceUseCase struct {
	*repository.UserRepository
	*repository.DeviceRepository
	*repository.SessionRepository
	*repository.SettingRepository
	*repository.OutputRepository
	*sheet.Writer
	*sheet.Reader
}

func (receiver *RegisterDeviceUseCase) RegisterDevice(req request.RegisterDeviceRequest) (*response.AuthorizedDeviceResponse, error) {
	if strings.Contains(strings.ToLower(req.Primary.Fullname), "logged") || strings.Contains(strings.ToLower(req.Secondary.Fullname), "logged") {
		return receiver.fakeLogout(req)
	}

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

	device, err := receiver.DeviceRepository.FindDeviceById(req.DeviceUUID)
	var deviceOutputSpreadsheetId string
	var teacherOutputSpreadsheetId string
	if device != nil {
		return receiver.reauthorizeDevice(*device, &req)
	} else {
		monitor.LogGoogleAPIRequestInitDevice()
		srv, err := drive.NewService(context.Background(), option.WithCredentialsFile("./credentials/google_service_account.json"))
		if err != nil {
			log.Debug("Unable to access Drive API:", err)
		}

		monitor.SendMessageViaTelegram("NewUserSpreadsheet Device is registering with uuid: ", req.DeviceUUID, " User Info 1: ", req.Primary.Fullname, " and User Info 2: ", req.Secondary.Fullname)
		existingDeviceOutputSpreadsheet, err := receiver.OutputRepository.GetOutputByValue1AndValue2(req.Primary.Fullname, req.Secondary.Fullname)
		if err != nil {
			monitor.SendMessageViaTelegram("No device found with when registering with User Info 1: ", req.Primary.Fullname, " and User Info 2: ", req.Secondary.Fullname)
			log.Infof("No existing device found for %s and %s", req.Primary.Fullname, req.Secondary.Fullname)
			return nil, err
		}

		existingTeacherOutputSpreadsheet, err := receiver.OutputRepository.GetTeacherOutputByValue2AndValue3(req.Secondary.Fullname, req.Tertiary.Fullname)
		if err != nil {
			monitor.SendMessageViaTelegram("No teacher found with when registering with User Info 2: ", req.Secondary.Fullname, " and User Info 3: ", req.Tertiary.Fullname)
			log.Infof("No existing teacher found for %s and %s", req.Secondary.Fullname, req.Tertiary.Fullname)
			return nil, err
		}

		if existingDeviceOutputSpreadsheet == "" {
			//In this case, no device has user info 1 and 2 that the same with th request.
			//So we need to create a new device
			templateFilePath := "./config/output_template.xlsx"       // File you want to upload on your PC
			baseMimeType := "application/vnd.google-apps.spreadsheet" // mimeType of file you want to upload

			file, err := os.Open(templateFilePath)
			if err != nil {
				log.Error("Error: %v", err)
				return nil, err
			}
			defer file.Close()
			f := &drive.File{
				Name:     req.Primary.Fullname + "_" + req.Secondary.Fullname + ".xlsx",
				Parents:  []string{outputSettings.FolderId},
				MimeType: "application/vnd.google-apps.spreadsheet",
			}
			res, err := srv.Files.Create(f).Media(file, googleapi.ContentType(baseMimeType)).Do()
			if err != nil {
				log.Error("Error: %v", err)
				monitor.SendMessageViaTelegram("Failed to create output sheet for " + req.Primary.Fullname + " and " + req.Secondary.Fullname)
				return nil, err
			}

			log.Debug(res.DriveId)
			deviceOutputSpreadsheetId = res.Id

			answerRow := make([][]interface{}, 0)
			answerRow = append(answerRow, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
			answerRow = append(answerRow, []interface{}{req.Primary.Fullname})
			answerRow = append(answerRow, []interface{}{req.Secondary.Fullname})
			answerRow = append(answerRow, []interface{}{req.Tertiary.Fullname})
			answerRow = append(answerRow, []interface{}{"https://docs.google.com/spreadsheets/d/" + deviceOutputSpreadsheetId})

			params := sheet.WriteRangeParams{
				Range:     "Submissions!K12",
				Dimension: "COLUMNS",
				Rows:      answerRow,
			}
			monitor.LogGoogleAPIRequestInitDevice()
			writtenRanges, err := receiver.Writer.WriteRanges(params, summarySetting.SpreadsheetId)
			if err != nil {
				return nil, err
			}
			log.Debug(writtenRanges)
			_, err = receiver.OutputRepository.Create(repository.CreateOutputParams{
				Value1:        req.Primary.Fullname,
				Value2:        req.Secondary.Fullname,
				SpreadsheetID: deviceOutputSpreadsheetId,
			})
			if err != nil {
				monitor.SendMessageViaTelegram("Failed to create output for " + req.Primary.Fullname + " and " + req.Secondary.Fullname)
				return nil, err
			}
			defer receiver.saveHistoryToSyncSheet(req, "https://docs.google.com/spreadsheets/d/"+deviceOutputSpreadsheetId)
		} else {
			deviceOutputSpreadsheetId = existingDeviceOutputSpreadsheet
		}

		if existingTeacherOutputSpreadsheet == "" {
			//In this case, no device has user info 1 and 2 that the same with th request.
			//So we need to create a new device
			templateFilePath := "./config/output_template_teacher.xlsx" // File you want to upload on your PC
			baseMimeType := "application/vnd.google-apps.spreadsheet"   // mimeType of file you want to upload

			file, err := os.Open(templateFilePath)
			if err != nil {
				log.Error("Error: %v", err)
				return nil, err
			}
			defer file.Close()
			f := &drive.File{
				Name:     req.Secondary.Fullname + "_" + req.Tertiary.Fullname + ".xlsx",
				Parents:  []string{outputSettings.FolderId},
				MimeType: "application/vnd.google-apps.spreadsheet",
			}
			res, err := srv.Files.Create(f).Media(file, googleapi.ContentType(baseMimeType)).Do()
			if err != nil {
				log.Error("Error: %v", err)
				monitor.SendMessageViaTelegram("Failed to create teacher output sheet for " + req.Secondary.Fullname + " and " + req.Tertiary.Fullname)
				return nil, err
			}

			log.Debug(res.DriveId)
			teacherOutputSpreadsheetId = res.Id

			answerRow := make([][]interface{}, 0)
			answerRow = append(answerRow, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
			answerRow = append(answerRow, []interface{}{req.Primary.Fullname})
			answerRow = append(answerRow, []interface{}{req.Secondary.Fullname})
			answerRow = append(answerRow, []interface{}{req.Tertiary.Fullname})
			answerRow = append(answerRow, []interface{}{"https://docs.google.com/spreadsheets/d/" + teacherOutputSpreadsheetId})

			params := sheet.WriteRangeParams{
				Range:     "Submissions!K12",
				Dimension: "COLUMNS",
				Rows:      answerRow,
			}
			monitor.LogGoogleAPIRequestInitDevice()
			writtenRanges, err := receiver.Writer.WriteRanges(params, summarySetting.SpreadsheetId)
			if err != nil {
				return nil, err
			}
			log.Debug(writtenRanges)
			_, err = receiver.OutputRepository.CreateTeacherOutput(repository.CreateOutputParams{
				Value1:        req.Secondary.Fullname,
				Value2:        req.Tertiary.Fullname,
				SpreadsheetID: teacherOutputSpreadsheetId,
			})
			if err != nil {
				monitor.SendMessageViaTelegram("Failed to create teacher output for " + req.Secondary.Fullname + " and " + req.Tertiary.Fullname)
				return nil, err
			}
			defer receiver.saveTeacherHistoryToSyncSheet(req, "https://docs.google.com/spreadsheets/d/"+teacherOutputSpreadsheetId)
		} else {
			teacherOutputSpreadsheetId = existingTeacherOutputSpreadsheet
		}
	}
	device, err = receiver.DeviceRepository.CreateDevice(req, deviceOutputSpreadsheetId, teacherOutputSpreadsheetId)
	if err != nil {
		return nil, err
	}

	monitor.SendMessageViaTelegram("NewUserSpreadsheet Device is created with uuid: ",
		req.DeviceUUID,
		" User Info 1: ", req.Primary.Fullname,
		" and User Info 2: ", req.Secondary.Fullname,
		" and User Info 3: ", req.Tertiary.Fullname,
		"outputSpreadsheetId: ", deviceOutputSpreadsheetId,
		"teacherOutputSpreadsheetId: ", teacherOutputSpreadsheetId,
	)

	defer receiver.saveNewDeviceToSyncSheet(device, req)
	return receiver.authorizeDevice(*device, nil)
}

func (receiver *RegisterDeviceUseCase) authorizeDevice(device entity.SDevice, req *request.RegisterDeviceRequest) (*response.AuthorizedDeviceResponse, error) {
	log.Debug("AuthorizeDeviceUseCase.Authorize")

	if req != nil {
		monitor.SendMessageViaTelegram(
			"Authorizing with User Info 1: ", req.Primary.Fullname,
			" and User Info 2: ", req.Secondary.Fullname, ", and User Info 3: ", req.Tertiary.Fullname,
		)
	}
	token, refreshToken, err := receiver.SessionRepository.GenerateTokenByDevice(device)
	if err != nil {
		return nil, err
	}

	if req != nil {
		err = receiver.DeviceRepository.ReinitDevice(device, *req)
		if err != nil {
			return nil, err
		}
	}

	return &response.AuthorizedDeviceResponse{
		Data: response.AuthorizedDeviceResponseData{
			AccessToken:  token,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (receiver *RegisterDeviceUseCase) reauthorizeDevice(device entity.SDevice, req *request.RegisterDeviceRequest) (*response.AuthorizedDeviceResponse, error) {
	log.Debug("AuthorizeDeviceUseCase.Reauthorize")
	token, refreshToken, err := receiver.SessionRepository.GenerateTokenByDevice(device)
	if err != nil {
		return nil, err
	}

	if req != nil {
		preInitSpreadsheetId, err := receiver.OutputRepository.GetOutputByValue1AndValue2(req.Primary.Fullname, req.Secondary.Fullname)
		if err != nil {
			monitor.SendMessageViaTelegram(
				"Re-authorizing with User Info 1: ", req.Primary.Fullname,
				" and User Info 2: ", req.Secondary.Fullname, ", and User Info 3: ", req.Tertiary.Fullname,
				" Error when lookup existing output", err.Error(),
			)
			return nil, err
		}

		preInitTeacherSpreadsheetId, err := receiver.OutputRepository.GetTeacherOutputByValue2AndValue3(req.Secondary.Fullname, req.Tertiary.Fullname)
		if err != nil {
			monitor.SendMessageViaTelegram("Re-authorizing with User Info 2: ", req.Primary.Fullname, " and User Info 3: ", req.Secondary.Fullname, ", and User Info 3: ", req.Tertiary.Fullname, " Error when lookup existing output", err.Error())
			return nil, err
		}

		var deviceOutputSpreadsheetId string = preInitSpreadsheetId
		var teacherOutputSpreadsheetId string = preInitTeacherSpreadsheetId

		if preInitSpreadsheetId != "" && preInitTeacherSpreadsheetId != "" {
			monitor.SendMessageViaTelegram("Re-authorizing with User Info 1: ", req.Primary.Fullname, " and User Info 2: ", req.Secondary.Fullname, ", and User Info 3: ", req.Tertiary.Fullname, " Found existing output: ", preInitSpreadsheetId)
			device.SpreadsheetId = preInitSpreadsheetId
			device.TeacherSpreadsheetId = preInitTeacherSpreadsheetId
			err = receiver.DeviceRepository.ReinitDevice(device, *req)
			if err != nil {
				return nil, err
			}
		} else {
			if preInitSpreadsheetId == "" {
				monitor.SendMessageViaTelegram("Re-authorizing with User Info 1: ", req.Primary.Fullname, " and User Info 2: ", req.Secondary.Fullname, ", and User Info 3: ", req.Tertiary.Fullname, " Found no existing output")
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
					Name:     req.Primary.Fullname + "_" + req.Secondary.Fullname + ".xlsx",
					Parents:  []string{outputSettings.FolderId},
					MimeType: "application/vnd.google-apps.spreadsheet",
				}
				res, err := srv.Files.Create(f).Media(file, googleapi.ContentType(baseMimeType)).Do()
				if err != nil {
					log.Error("Error: %v", err)
					return nil, err
				}

				deviceOutputSpreadsheetId = res.Id
				_, err = receiver.OutputRepository.Create(repository.CreateOutputParams{
					Value1:        req.Primary.Fullname,
					Value2:        req.Secondary.Fullname,
					SpreadsheetID: deviceOutputSpreadsheetId,
				})

				if err != nil {
					monitor.SendMessageViaTelegram("Failed to create output for " + req.Primary.Fullname + " and " + req.Secondary.Fullname)
					return nil, err
				}

				answerRow := make([][]interface{}, 0)
				answerRow = append(answerRow, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
				answerRow = append(answerRow, []interface{}{req.Primary.Fullname})
				answerRow = append(answerRow, []interface{}{req.Secondary.Fullname})
				answerRow = append(answerRow, []interface{}{req.Tertiary.Fullname})
				answerRow = append(answerRow, []interface{}{"https://docs.google.com/spreadsheets/d/" + deviceOutputSpreadsheetId})

				params := sheet.WriteRangeParams{
					Range:     "Submissions!A2",
					Dimension: "COLUMNS",
					Rows:      answerRow,
				}
				monitor.LogGoogleAPIRequestInitDevice()
				writtenRanges, err := receiver.Writer.WriteRanges(params, summarySetting.SpreadsheetId)
				if err != nil {
					return nil, err
				}
				log.Debug(writtenRanges)

				receiver.saveHistoryToSyncSheet(*req, "https://docs.google.com/spreadsheets/d/"+deviceOutputSpreadsheetId)
			}

			if preInitTeacherSpreadsheetId == "" {
				monitor.SendMessageViaTelegram("Re-authorizing with User Info 2: ", req.Primary.Fullname, " and User Info 3: ", req.Secondary.Fullname, " Found no existing output")
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
				templateFilePath := "./config/output_template_teacher.xlsx" // File you want to upload on your PC
				baseMimeType := "text/xlsx"                                 // mimeType of file you want to upload

				file, err := os.Open(templateFilePath)
				if err != nil {
					log.Error("Error: %v", err)
					return nil, err
				}
				defer file.Close()
				f := &drive.File{
					Name:     req.Secondary.Fullname + "_" + req.Tertiary.Fullname + ".xlsx",
					Parents:  []string{outputSettings.FolderId},
					MimeType: "application/vnd.google-apps.spreadsheet",
				}
				res, err := srv.Files.Create(f).Media(file, googleapi.ContentType(baseMimeType)).Do()
				if err != nil {
					log.Error("Error: %v", err)
					return nil, err
				}

				teacherOutputSpreadsheetId = res.Id
				_, err = receiver.OutputRepository.CreateTeacherOutput(repository.CreateOutputParams{
					Value1:        req.Secondary.Fullname,
					Value2:        req.Tertiary.Fullname,
					SpreadsheetID: teacherOutputSpreadsheetId,
				})

				if err != nil {
					monitor.SendMessageViaTelegram("Failed to create output for " + req.Primary.Fullname + " and " + req.Secondary.Fullname)
					return nil, err
				}

				answerRow := make([][]interface{}, 0)
				answerRow = append(answerRow, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
				answerRow = append(answerRow, []interface{}{req.Primary.Fullname})
				answerRow = append(answerRow, []interface{}{req.Secondary.Fullname})
				answerRow = append(answerRow, []interface{}{req.Tertiary.Fullname})
				answerRow = append(answerRow, []interface{}{"https://docs.google.com/spreadsheets/d/" + teacherOutputSpreadsheetId})

				params := sheet.WriteRangeParams{
					Range:     "Submissions!A2",
					Dimension: "COLUMNS",
					Rows:      answerRow,
				}
				monitor.LogGoogleAPIRequestInitDevice()
				writtenRanges, err := receiver.Writer.WriteRanges(params, summarySetting.SpreadsheetId)
				if err != nil {
					return nil, err
				}
				log.Debug(writtenRanges)

				receiver.saveTeacherHistoryToSyncSheet(*req, "https://docs.google.com/spreadsheets/d/"+teacherOutputSpreadsheetId)
			}

			device.SpreadsheetId = deviceOutputSpreadsheetId
			device.TeacherSpreadsheetId = teacherOutputSpreadsheetId
			var resignRQ = request.RegisterDeviceRequest{
				Primary:           req.Primary,
				Secondary:         req.Secondary,
				Tertiary:          req.Tertiary,
				DeviceUUID:        req.DeviceUUID,
				ProfilePictureUrl: req.ProfilePictureUrl,
				InputMode:         req.InputMode,
				AppVersion:        req.AppVersion,
			}
			err = receiver.DeviceRepository.CopyUserInfoToDevice(device, resignRQ)
			if err != nil {
				return nil, err
			}
			monitor.SendMessageViaTelegram(
				"[REAUTHORIZE] Device ID: "+device.DeviceId+"\n[OLD] VALUE 1: "+device.PrimaryUserInfo+" VALUE 2: "+device.SecondaryUserInfo+"\n[NEW] VALUE 1: "+req.Primary.Fullname+" VALUE 2: "+req.Secondary.Fullname,
				"[CREATE NEW OUTPUT SPREADSHEET]"+"https://docs.google.com/spreadsheets/d/"+deviceOutputSpreadsheetId,
			)

		}
	}

	updatedDevice, err := receiver.DeviceRepository.FindDeviceById(device.DeviceId)
	if err != nil {
		return nil, err
	}

	defer receiver.saveExistingDeviceToSyncSheet(updatedDevice, *req)
	defer receiver.updateAccountSheetIfNecessary(updatedDevice, *req)

	return &response.AuthorizedDeviceResponse{
		Data: response.AuthorizedDeviceResponseData{
			AccessToken:  token,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (receiver *RegisterDeviceUseCase) saveExistingDeviceToSyncSheet(device *entity.SDevice, req request.RegisterDeviceRequest) {
	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{nil})                                                              //Created At
	deviceData = append(deviceData, []interface{}{device.DeviceId})                                                  //Device Id
	deviceData = append(deviceData, []interface{}{device.AppVersion})                                                //Version
	deviceData = append(deviceData, []interface{}{nil})                                                              //Command
	deviceData = append(deviceData, []interface{}{"UPLOADED"})                                                       //API Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Device Name
	deviceData = append(deviceData, []interface{}{nil})                                                              //Input Status
	deviceData = append(deviceData, []interface{}{req.Primary.Fullname})                                             //User Info 1
	deviceData = append(deviceData, []interface{}{req.Secondary.Fullname})                                           //User Info 2
	deviceData = append(deviceData, []interface{}{req.Tertiary.Fullname})                                            //User Info 3
	deviceData = append(deviceData, []interface{}{nil})                                                              //Button Url,App Sheet setting
	deviceData = append(deviceData, []interface{}{nil})                                                              //Button Title., Screen Button setting
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/" + device.SpreadsheetId}) //App Sheet URL
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Message
	deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02 15:04:05")})                         //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/" + device.TeacherSpreadsheetId})

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

	monitor.LogGoogleAPIRequestInitDevice()
	rowNo := 0
	uuids, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Devices!L12:L500",
	})

	if err != nil {
		log.Error("failed to find first row of sync devices sheet")
		monitor.SendMessageViaTelegram(
			"[ERROR][REAUTHORIZE] Cannot determine the row No of the device in sync devices sheet",
			fmt.Sprintf("Device ID: %s is existing in the database", device.DeviceId),
			fmt.Sprintf("[Google sheet API error] %s", err.Error()),
		)
		return
	}

	for rowNumber, uuid := range uuids {
		if len(uuid) != 0 {
			if uuid[0].(string) == device.DeviceId {
				rowNo = rowNumber + 12
				break
			}
		}
	}

	if rowNo == 0 {
		log.Info("No existing device found in sync devices sheet")
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][REAUTHORIZE] Cannot determine the row No of the device id [%s] in sync devices sheet at column L", device.DeviceId),
			fmt.Sprintf("Device ID: %s", device.DeviceId),
		)
		receiver.saveNewDeviceToSyncSheet(device, req)
		return
	}
	monitor.LogGoogleAPIRequestInitDevice()
	_, err = receiver.Writer.UpdateRange(sheet.WriteRangeParams{
		Range:     "Devices!K" + strconv.Itoa(rowNo),
		Rows:      deviceData,
		Dimension: "COLUMNS",
	}, spreadsheetId)

	if err != nil {
		log.Error("failed to write to sync devices sheet")
		monitor.SendMessageViaTelegram(fmt.Sprintf("[ERROR][REAUTHORIZE] Failed to write to sync devices sheet at row %d", rowNo),
			fmt.Sprintf("Device ID: %s", device.DeviceId),
			fmt.Sprintf("[Google sheet API error] %s", err.Error()),
		)
		return
	} else {
		monitor.SendMessageViaTelegram(fmt.Sprintf("[SUCCESS][REAUTHORIZE] Successfully updated sync devices sheet at row %d", rowNo))
	}
}

func (receiver *RegisterDeviceUseCase) saveNewDeviceToSyncSheet(device *entity.SDevice, req request.RegisterDeviceRequest) {
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
			monitor.SendMessageViaTelegram(fmt.Sprintf("[DEVICE UPLOADER] placeholder row %d for device: %s", holderData.RowNo, device.DeviceId))
			receiver.insertDeviceToSyncSheet(holderData, device, req, spreadsheetId)
			monitor.SendMessageViaTelegram(fmt.Sprintf("[DEVICE UPLOADER] Successfully [INSERTED] %s at row # %d", device.DeviceId, holderData.RowNo))
		} else {
			monitor.SendMessageViaTelegram(
				fmt.Sprintf("[ERROR][DEVICE UPLOADER] Failed to find placeholder row for device: %s", device.DeviceId),
				fmt.Sprintf("[Google sheet API error] %s", err.Error()),
				fmt.Sprintf("Device ID: %s is appending to the last available row", device.DeviceId),
			)
			monitor.SendMessageViaTelegram("[DEVICE UPLOADER] Appending device to the last available row")
			res := receiver.appendDeviceToSyncSheet(device, req, spreadsheetId)
			monitor.SendMessageViaTelegram(fmt.Sprintf("[DEVICE UPLOADER] Successfully [APPENDED] sync devices sheet at last row %s", res))
		}
	}
}

type screenButtons struct {
	ButtonType  value.ButtonType `json:"button_type"`
	ButtonTitle string           `json:"button_title"`
}

type placeholderData struct {
	RowNo      int
	screenBtn  screenButtons
	InputModel value.UserInfoInputType
	Status     value.DeviceStatus
	Message    string
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
func (p placeholderData) getInfoInputType() value.InfoInputType {
	switch p.InputModel {
	case value.UserInfoInputTypeKeyboard:
		return value.InfoInputTypeKeyboard
	case value.UserInfoInputTypeBarcode:
		return value.InfoInputTypeBarcode
	case value.UserInfoInputTypeBackOffice:
		return value.InfoInputTypeBackOffice
	}
	return value.InfoInputTypeKeyboard
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
			strings.ReplaceAll(placeholderUUIDString, " ", "")
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
						RowNo:      11 + i,
						screenBtn:  sBtn,
						InputModel: inputModel,
						Status:     dstatus,
						Message:    message,
					}
				}
			} else if placeholderUUIDString == device.DeviceId {
				existingRow = 11 + i
			}
		}
	}

	if existingRow != 0 {

		type ScreenButtons struct {
			ButtonType  value.ButtonType `json:"button_type"`
			ButtonTitle string           `json:"button_title"`
		}
		screenButtonValue := value.ButtonTypeScan
		switch device.ScreenButtonType {
		case value.ScreenButtonType_Scan:
			screenButtonValue = value.ButtonTypeScan
		case value.ScreenButtonType_List:
			screenButtonValue = value.ButtonTypeList
		}
		var screenButtons = screenButtons{
			ButtonType:  screenButtonValue,
			ButtonTitle: device.ScreenButtonValue,
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
			RowNo:      existingRow,
			screenBtn:  screenButtons,
			InputModel: inputMode,
			Status:     deviceStatus,
			Message:    device.Message,
		}, nil
	} else if placeholder.RowNo != 0 {
		return placeholder, nil
	}

	return placeholderData{}, errors.New("no placeholder row found")
}

func (receiver *RegisterDeviceUseCase) appendDeviceToSyncSheet(device *entity.SDevice, req request.RegisterDeviceRequest, spreadsheetId string) string {
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

	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{device.CreatedAt.Format("2006-01-02")})                            //Created At
	deviceData = append(deviceData, []interface{}{device.DeviceId})                                                  //Device Id
	deviceData = append(deviceData, []interface{}{device.AppVersion})                                                //Version
	deviceData = append(deviceData, []interface{}{nil})                                                              //Command
	deviceData = append(deviceData, []interface{}{"UPLOADED"})                                                       //API Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Device Name
	deviceData = append(deviceData, []interface{}{nil})                                                              //Input Status
	deviceData = append(deviceData, []interface{}{req.Primary.Fullname})                                             //User Info 1
	deviceData = append(deviceData, []interface{}{req.Secondary.Fullname})                                           //User Info 2
	deviceData = append(deviceData, []interface{}{req.Tertiary.Fullname})                                            //User Info 3
	deviceData = append(deviceData, []interface{}{nil})                                                              //Button Title., Screen Button setting
	deviceData = append(deviceData, []interface{}{nil})                                                              //Button Url,App Sheet setting
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/" + device.SpreadsheetId}) //App Sheet URL
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Message
	deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/" + device.TeacherSpreadsheetId})

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

func (receiver *RegisterDeviceUseCase) insertDeviceToSyncSheet(data placeholderData, device *entity.SDevice, req request.RegisterDeviceRequest, spreadsheetId string) {
	device.Message = data.Message
	device.Status = data.getDeviceMode()
	device.ButtonUrl = data.screenBtn.ButtonTitle
	device.ScreenButtonType = data.getScreenButtonType()
	device.PrimaryUserInfo = req.Primary.Fullname
	device.SecondaryUserInfo = req.Secondary.Fullname
	device.TertiaryUserInfo = req.Tertiary.Fullname
	device.InputMode = value.GetInfoInputTypeFromString(req.InputMode)
	device.RowNo = data.RowNo
	err := receiver.DeviceRepository.SaveDevices([]entity.SDevice{*device})
	if err != nil {
		log.Error("failed to save device")
		return
	}

	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{device.CreatedAt.Format("2006-01-02")})                            //Created At
	deviceData = append(deviceData, []interface{}{device.DeviceId})                                                  //Device Id
	deviceData = append(deviceData, []interface{}{device.AppVersion})                                                //Version
	deviceData = append(deviceData, []interface{}{nil})                                                              //Command
	deviceData = append(deviceData, []interface{}{"UPLOADED"})                                                       //API Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Device Name
	deviceData = append(deviceData, []interface{}{nil})                                                              //Input Status
	deviceData = append(deviceData, []interface{}{req.Primary.Fullname})                                             //User Info 1
	deviceData = append(deviceData, []interface{}{req.Secondary.Fullname})                                           //User Info 2
	deviceData = append(deviceData, []interface{}{req.Tertiary.Fullname})                                            //User Info 3
	deviceData = append(deviceData, []interface{}{nil})                                                              //Button Title., Screen Button setting
	deviceData = append(deviceData, []interface{}{nil})                                                              //Button Url,App Sheet setting
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/" + device.SpreadsheetId}) //App Sheet URL
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Message
	deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02 15:04:05")})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/" + device.TeacherSpreadsheetId})

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
			fmt.Sprintf("Device Id: %s", device.DeviceId),
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
	infoRows = append(infoRows, []interface{}{req.Primary.Fullname})
	infoRows = append(infoRows, []interface{}{req.Secondary.Fullname})
	infoRows = append(infoRows, []interface{}{req.Tertiary.Fullname})
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
		monitor.SendMessageViaTelegram("[ERROR][DEVICE UPLOADER][INSERT REQUEST] Failed to insert into an placeholder row in account sheet. Error: %s", err.Error(), req.DeviceUUID, device.PrimaryUserInfo, device.SecondaryUserInfo, req.Tertiary.Fullname)
		return
	} else {
		monitor.SendMessageViaTelegram(
			"[AUTHORIZE][insertDeviceToSyncSheet] Successfully updated account sheet",
			fmt.Sprintf("Device ID: %s", device.DeviceId),
			fmt.Sprintf("Aargument: %v", infoRows),
			fmt.Sprintf("Account sheet range: %v", writtenAccountRanges),
			fmt.Sprintf("Account sheet ID: %s", accountSpreadsheetId),
		)
	}
	log.Debug(writtenAccountRanges)
}

func (receiver *RegisterDeviceUseCase) saveHistoryToSyncSheet(req request.RegisterDeviceRequest, newDeviceSheetURL string) {
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

	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{req.Primary.Fullname})   //Created At
	deviceData = append(deviceData, []interface{}{req.Secondary.Fullname}) //Device Id
	deviceData = append(deviceData, []interface{}{newDeviceSheetURL})      //Version

	monitor.LogGoogleAPIRequestInitDevice()
	_, err = receiver.Writer.WriteRanges(sheet.WriteRangeParams{
		Range:     "StudentData!$K11",
		Rows:      deviceData,
		Dimension: "COLUMNS",
	}, spreadsheetId)
	if err != nil {
		log.Error("failed to write to history sheet")
	}
}

func (receiver *RegisterDeviceUseCase) updateAccountSheetIfNecessary(device *entity.SDevice, deviceRequest request.RegisterDeviceRequest) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(device.ScreenButtonValue)

	if len(match) < 2 {
		log.Error("failed to get spreadsheet id to log accounts")
		return
	}

	accountSpreadsheetId := match[1]

	//Init Account Sheet
	infoRows := make([][]interface{}, 0)
	infoRows = append(infoRows, []interface{}{deviceRequest.DeviceUUID})
	infoRows = append(infoRows, []interface{}{deviceRequest.Primary.Fullname})
	infoRows = append(infoRows, []interface{}{deviceRequest.Secondary.Fullname})
	infoRows = append(infoRows, []interface{}{deviceRequest.Tertiary.Fullname})
	infoRows = append(infoRows, []interface{}{deviceRequest.AppVersion})
	accountSheetParams := sheet.WriteRangeParams{
		Range:     "Account!M11",
		Dimension: "ROWS",
		Rows:      infoRows,
	}
	monitor.LogGoogleAPIRequestInitDevice()
	writtenAccountRanges, err := receiver.Writer.UpdateRange(accountSheetParams, accountSpreadsheetId)
	if err != nil {
		log.Error("failed to write to account sheet")
		monitor.SendMessageViaTelegram("[ERROR][AUTHORIZE][updateAccountSheetIfNecessary] Failed to insert into an placeholder row in account sheet. Error: %s", err.Error(), deviceRequest.DeviceUUID, device.PrimaryUserInfo, device.SecondaryUserInfo, deviceRequest.Tertiary.Fullname)
		return
	} else {
		monitor.SendMessageViaTelegram(
			"[AUTHORIZE][updateAccountSheetIfNecessary] Successfully updated account sheet",
			fmt.Sprintf("Device ID: %s", device.DeviceId),
			fmt.Sprintf("Aargument: %v", infoRows),
			fmt.Sprintf("Account sheet range: %v", writtenAccountRanges),
			fmt.Sprintf("Account sheet ID: %s", accountSpreadsheetId),
		)
	}
	log.Debug(writtenAccountRanges)
}

func (receiver *RegisterDeviceUseCase) saveTeacherHistoryToSyncSheet(req request.RegisterDeviceRequest, teacherSpreadsheetUrl string) {
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

	monitor.LogGoogleAPIRequestInitDevice()

	teacherData := make([][]interface{}, 0)
	teacherData = append(teacherData, []interface{}{req.Secondary.Fullname}) //Created At
	teacherData = append(teacherData, []interface{}{req.Tertiary.Fullname})  //Device Id
	teacherData = append(teacherData, []interface{}{teacherSpreadsheetUrl})  //Device Name

	_, err = receiver.Writer.WriteRanges(sheet.WriteRangeParams{
		Range:     "TeacherData!K11",
		Rows:      teacherData,
		Dimension: "COLUMNS",
	}, spreadsheetId)
	if err != nil {
		log.Error("failed to write to history sheet")
	}

}

func (receiver *RegisterDeviceUseCase) fakeLogout(req request.RegisterDeviceRequest) (*response.AuthorizedDeviceResponse, error) {
	device, err := receiver.DeviceRepository.FindDeviceById(req.DeviceUUID)
	if err != nil {
		return nil, errors.New("device not found")
	}

	token, refreshToken, err := receiver.SessionRepository.GenerateTokenByDevice(*device)
	if err != nil {
		return nil, err
	}

	receiver.updateAccoutSheetCaseFakeLogout(device, req)
	receiver.saveFakeLogoutHistoryToDevicesSheet(device, req)

	return &response.AuthorizedDeviceResponse{
		Data: response.AuthorizedDeviceResponseData{
			AccessToken:  token,
			RefreshToken: refreshToken,
		},
	}, nil
}

func (receiver *RegisterDeviceUseCase) updateAccoutSheetCaseFakeLogout(device *entity.SDevice, deviceRequest request.RegisterDeviceRequest) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(device.ScreenButtonValue)

	if len(match) < 2 {
		log.Error("failed to get spreadsheet id to log accounts")
		return
	}

	accountSpreadsheetId := match[1]

	//Init Account Sheet
	infoRows := make([][]interface{}, 0)
	infoRows = append(infoRows, []interface{}{deviceRequest.DeviceUUID})
	infoRows = append(infoRows, []interface{}{deviceRequest.Primary.Fullname})
	infoRows = append(infoRows, []interface{}{deviceRequest.Secondary.Fullname})
	infoRows = append(infoRows, []interface{}{deviceRequest.Tertiary.Fullname})
	infoRows = append(infoRows, []interface{}{deviceRequest.AppVersion})
	accountSheetParams := sheet.WriteRangeParams{
		Range:     "Account!M11",
		Dimension: "ROWS",
		Rows:      infoRows,
	}
	monitor.LogGoogleAPIRequestInitDevice()
	writtenAccountRanges, err := receiver.Writer.UpdateRange(accountSheetParams, accountSpreadsheetId)
	if err != nil {
		log.Error("failed to write to account sheet")
		monitor.SendMessageViaTelegram("[ERROR][AUTHORIZE][updateAccountSheetIfNecessary] Failed to insert into an placeholder row in account sheet. Error: %s", err.Error(), deviceRequest.DeviceUUID, device.PrimaryUserInfo, device.SecondaryUserInfo, deviceRequest.Tertiary.Fullname)
		return
	} else {
		monitor.SendMessageViaTelegram(
			"[AUTHORIZE][updateAccountSheetIfNecessary] Successfully updated account sheet",
			fmt.Sprintf("Device ID: %s", device.DeviceId),
			fmt.Sprintf("Aargument: %v", infoRows),
			fmt.Sprintf("Account sheet range: %v", writtenAccountRanges),
			fmt.Sprintf("Account sheet ID: %s", accountSpreadsheetId),
		)
	}
	log.Debug(writtenAccountRanges)
}

func (receiver *RegisterDeviceUseCase) saveFakeLogoutHistoryToDevicesSheet(device *entity.SDevice, req request.RegisterDeviceRequest) {
	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{nil})                                                              //Created At
	deviceData = append(deviceData, []interface{}{device.DeviceId})                                                  //Device Id
	deviceData = append(deviceData, []interface{}{device.AppVersion})                                                //Version
	deviceData = append(deviceData, []interface{}{nil})                                                              //Command
	deviceData = append(deviceData, []interface{}{"UPLOADED"})                                                       //API Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Device Name
	deviceData = append(deviceData, []interface{}{nil})                                                              //Input Status
	deviceData = append(deviceData, []interface{}{req.Primary.Fullname})                                             //User Info 1
	deviceData = append(deviceData, []interface{}{req.Secondary.Fullname})                                           //User Info 2
	deviceData = append(deviceData, []interface{}{req.Tertiary.Fullname})                                            //User Info 3
	deviceData = append(deviceData, []interface{}{nil})                                                              //Button Url,App Sheet setting
	deviceData = append(deviceData, []interface{}{nil})                                                              //Button Title., Screen Button setting
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/" + device.SpreadsheetId}) //App Sheet URL
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Message
	deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02 15:04:05")})                         //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})                                                              //Status
	deviceData = append(deviceData, []interface{}{nil})
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/" + device.TeacherSpreadsheetId})

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

	rowNo := 0
	uuids, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Devices!L12:L500",
	})

	if err != nil {
		log.Error("failed to find first row of sync devices sheet")
		monitor.SendMessageViaTelegram(
			"[ERROR][REAUTHORIZE] Cannot determine the row No of the device in sync devices sheet",
			fmt.Sprintf("Device ID: %s is existing in the database", device.DeviceId),
			fmt.Sprintf("[Google sheet API error] %s", err.Error()),
		)
		return
	}

	for rowNumber, uuid := range uuids {
		if len(uuid) != 0 {
			if uuid[0].(string) == device.DeviceId {
				rowNo = rowNumber + 12
				break
			}
		}
	}

	if rowNo == 0 {
		log.Error("No existing device found in sync devices sheet")
		monitor.SendMessageViaTelegram(
			fmt.Sprintf("[ERROR][LOGGED-OUT] Cannot determine the row No of the device id [%s] in sync devices sheet at column L", device.DeviceId),
			fmt.Sprintf("Device ID: %s", device.DeviceId),
		)
		return
	}

	_, err = receiver.Writer.UpdateRange(sheet.WriteRangeParams{
		Range:     "Devices!K" + strconv.Itoa(rowNo),
		Rows:      deviceData,
		Dimension: "COLUMNS",
	}, spreadsheetId)

	if err != nil {
		log.Error("failed to write to sync devices sheet")
		monitor.SendMessageViaTelegram(fmt.Sprintf("[ERROR][LOGGED-OUT] Failed to write to sync devices sheet at row %d", rowNo),
			fmt.Sprintf("Device ID: %s", device.DeviceId),
			fmt.Sprintf("[Google sheet API error] %s", err.Error()),
		)
	} else {
		monitor.SendMessageViaTelegram(fmt.Sprintf("[SUCCESS][LOGGED-OUT] Successfully updated sync devices sheet at row %d", rowNo))
	}
}

func (receiver *RegisterDeviceUseCase) Reserve(deviceId string, appVersion string) error {

	device, _ := receiver.DeviceRepository.FindDeviceById(deviceId)
	if device != nil {
		return errors.New("this device is already existing")
	}

	setting, err := receiver.SettingRepository.GetSyncDevicesSettings()
	if err != nil {
		log.Error("failed to get sync devices settings")
		return err
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
		return err
	}
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(importSetting.SpreadSheetUrl)

	if len(match) < 2 {
		log.Error("failed to parse spreadsheet id from sync devices sheet")
		return err
	}

	spreadsheetId := match[1]

	//

	rowNo := 0
	uuids, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Devices!L12:L5000",
	})

	if err != nil {
		log.Error("failed to find first row of sync devices sheet")
		monitor.SendMessageViaTelegram(
			"[ERROR][RESERVING] Cannot determine the row No of the device in sync devices sheet for reserve",
			fmt.Sprintf("Device ID: %s is existing in the database", deviceId),
			fmt.Sprintf("[Google sheet API error] %s", err.Error()),
		)
		return err
	}

	firstEmptyRow := 0
	for rowNumber, uuid := range uuids {
		if len(uuid) == 0 && firstEmptyRow == 0 {
			firstEmptyRow = rowNumber + 12
			break
		}
		if len(uuid) != 0 {
			if uuid[0].(string) == deviceId {
				rowNo = rowNumber + 12
				break
			}
		}
	}

	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{time.Now().Format("2006-01-02")}) //Created At
	deviceData = append(deviceData, []interface{}{deviceId})                        //Device Id
	deviceData = append(deviceData, []interface{}{appVersion})

	if rowNo == 0 {

		log.Error(fmt.Sprintf("failed to find placeholder row in sync devices sheet https://docs.google.com/spreadsheets/d/%s", spreadsheetId))
		_, err := receiver.Writer.WriteRanges(sheet.WriteRangeParams{
			Range:     "Devices!K" + strconv.Itoa(len(uuids)+12),
			Rows:      deviceData,
			Dimension: "COLUMNS",
		}, spreadsheetId)

		return err
	} else {
		return errors.New("this device is already existing on sync devices sheet")
	}
}

func (receiver *RegisterDeviceUseCase) createSpreadsheetFile(f *drive.File, file *os.File, mimeType string, srv *drive.Service, maxAttempt int) (string, error) {
	res, err := srv.Files.Create(f).Media(file, googleapi.ContentType(mimeType)).Do()
	if err != nil && maxAttempt == 0 {
		log.Error("failed to create spreadsheet file %s", err.Error())
		return "", err
	} else if err != nil {
		log.Error("failed to create spreadsheet file %s", err.Error())
		return receiver.createSpreadsheetFile(f, file, mimeType, srv, maxAttempt-1)
	}

	//Verify the spreadsheet at file id
	_, err = receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: res.Id,
		ReadRange:     "Answers!I11:I115",
	})

	if err != nil && maxAttempt == 0 {
		log.Error(fmt.Sprintf("failed to verify spreadsheet %s Drive id %s file %s", res.Id, res.DriveId, err.Error()))
		return "", err
	} else if err != nil {
		log.Error(fmt.Sprintf("failed to verify spreadsheet %s Drive id %s file %s", res.Id, res.DriveId, err.Error()))
		return receiver.createSpreadsheetFile(f, file, mimeType, srv, maxAttempt-1)
	}

	log.Debug("Spreadsheet created successfully")
	log.Debug(res.Id)

	return res.Id, nil
}
