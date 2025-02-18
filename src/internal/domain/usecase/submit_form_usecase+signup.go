package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/parameters"
	"sen-global-api/internal/domain/request"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

func (receiver *SubmitFormUseCase) SubmitSignUpForm(form entity.SForm, rq request.SubmitFormRequest) error {
	// setting, err := receiver.SettingRepository.GetRegistrationSubmissionSetting()
	// if err != nil {
	// 	return err
	// }

	// if setting == nil {
	// 	return errors.New("registration submission setting is not set")
	// }

	// type summarySetting struct {
	// 	SpreadsheetId string `json:"spreadsheet_id"`
	// }

	// var summary summarySetting
	// err = json.Unmarshal(setting.Settings, &summary)
	// if err != nil {
	// 	return err
	// }

	// re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	// match := re.FindStringSubmatch(summary.SpreadsheetId)

	// if len(match) < 2 {
	// 	return errors.New("invalid spreadsheet url from sign up submission setting")
	// }

	// spreadsheetId := match[1]

	submissionItems := make([]repository.SubmissionDataItem, 0)
	questions, err := receiver.QuestionRepository.GetQuestionsByIDs(Map(rq.Answers, func(answer request.Answer) string { return answer.QuestionId }))
	if err != nil {
		return fmt.Errorf("system cannot find questions for this form: %s", form.Name)
	}

	for _, answer := range rq.Answers {
		for _, question := range questions {
			if answer.QuestionId == question.QuestionId {
				var msg *repository.Messaging = nil
				if answer.Messaging != nil {
					msg = &repository.Messaging{
						Email:        answer.Messaging.Email,
						Value3:       answer.Messaging.Value3,
						MessageBox:   answer.Messaging.MessageBox,
						QuestionType: answer.Messaging.QuestionType,
					}
				}
				submissionItems = append(submissionItems, repository.SubmissionDataItem{
					QuestionId: question.QuestionId,
					Question:   question.Question,
					Answer:     answer.Answer,
					Messaging:  msg,
				})
			}
		}
	}
	submissionData := repository.SubmissionData{
		Items: submissionItems,
	}

	createSubmissionParmas := repository.CreateSubmissionParams{
		FormId:         form.ID,
		DeviceId:       "SignedUp_At_" + time.Now().Format("20060102150405"),
		SubmissionData: submissionData,
		OpenedAt:       rq.OpenedAt,
	}
	err = receiver.SubmissionRepository.CreateSubmission(createSubmissionParmas)
	if err != nil {
		log.Error("SubmitFormUseCase.SubmitSignUpForm", err)
		return errors.New("system cannot handle the submission")
	}

	return nil
}

func (receiver *SubmitFormUseCase) SubmitSignUpMemoryForm(form entity.SForm, rq request.SubmitFormRequest) error {
	if rq.DeviceId == "" {
		return errors.New("device id is required")
	}

	err := receiver.SubmitSignUpForm(form, rq)
	if err != nil {
		return err
	}

	err = receiver.createSignUpMemoryForm(form, rq)
	if err != nil {
		return err
	}

	return nil
}

func (receiver *SubmitFormUseCase) createSignUpMemoryForm(form entity.SForm, rq request.SubmitFormRequest) error {
	outputSettingsData, err := receiver.SettingRepository.GetOutputSettings()
	if err != nil {
		return err
	}

	var outputSettings OutputSetting
	if outputSettingsData != nil {
		err = json.Unmarshal([]byte(outputSettingsData.Settings), &outputSettings)
		if err != nil {
			return err
		}
	}

	registrationFormSetting, err := receiver.SettingRepository.GetRegistrationFormSetting()
	if err != nil {
		return err
	}

	type summarySetting struct {
		SpreadsheetId string `json:"spreadsheet_id"`
	}

	var summary summarySetting
	err = json.Unmarshal(registrationFormSetting.Settings, &summary)
	if err != nil {
		return err
	}

	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(summary.SpreadsheetId)

	if len(match) < 2 {
		return errors.New("invalid spreadsheet url from sign up submission setting")
	}

	signUpFormSpreadsheetId := match[1]

	//Clone sign Up Registration Form into memory form

	targetSpreadsheetName := "SENBOX.ORG/SIGN-UP[Memory-Form][Device-" + rq.DeviceId + "]"
	targetSheetName := "Sign Up"
	duplicateSpreadsheetResult, err := receiver.Writer.DuplicateSpreadsheet(sheet.DuplicateSpreadsheetParams{
		SourceSpreadsheetId:   signUpFormSpreadsheetId,
		TargetSpreadsheetName: targetSpreadsheetName,
		TargetSheetName:       targetSheetName,
	})

	if err != nil {
		return err
	}

	log.Debug("duplicateSpreadsheetResult ", duplicateSpreadsheetResult)
	// Retrieve the existing file to get the current parents
	file, err := receiver.DriveService.Files.Get(duplicateSpreadsheetResult.SpreadsheetId).Fields("parents").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve file: %v", err)
	}

	previousParents := ""
	if len(file.Parents) > 0 {
		previousParents = file.Parents[0]
	}

	updateResult, err := receiver.DriveService.Files.Update(duplicateSpreadsheetResult.SpreadsheetId, nil).
		AddParents(outputSettings.FolderId).
		RemoveParents(previousParents).
		Fields("id, parents").
		Do()
	if err != nil {
		log.Fatalf("Unable to move file: %v", err)
		return err
	}

	log.Debug("updateResult ", updateResult)
	//Create new form to db
	formQuestions, err := receiver.FormQuestionRepository.GetFormQuestionsByForm(form)

	if err != nil {
		return err
	}

	questions, err := receiver.QuestionRepository.GetQuestionsByIDs(Map(formQuestions, func(question entity.SFormQuestion) string { return question.QuestionId }))
	if err != nil {
		return err
	}

	rawQuestions := make([]entity.SQuestion, 0)
	for _, question := range questions {

		questionIndex := strings.Replace(question.QuestionId, fmt.Sprintf("%s_%s_", strings.ToUpper(form.Note), form.SpreadsheetId), "", -1)
		questionID := fmt.Sprintf("%s_%s_%s", strings.ToUpper(targetSpreadsheetName), duplicateSpreadsheetResult.SpreadsheetId, questionIndex)

		rawQuestions = append(rawQuestions, entity.SQuestion{
			QuestionId:     questionID,
			QuestionName:   question.QuestionName,
			QuestionType:   question.QuestionType,
			Question:       question.Question,
			Attributes:     question.Attributes,
			Status:         question.Status,
			Set:            question.Set,
			EnableOnMobile: question.EnableOnMobile,
		})
	}

	// setting, err := receiver.SettingRepository.GetRegistrationSubmissionSetting()
	// if err != nil {
	// 	return err
	// }

	// if setting == nil {
	// 	return errors.New("registration submission setting is not set")
	// }

	// err = json.Unmarshal(setting.Settings, &summary)
	// if err != nil {
	// 	return err
	// }

	// re = regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	// match = re.FindStringSubmatch(summary.SpreadsheetId)

	// if len(match) < 2 {
	// 	return errors.New("invalid spreadsheet url from sign up submission setting")
	// }

	// submissionSpreadsheetID := match[1]

	newForm, err := receiver.FormRepository.SaveForm(parameters.SaveFormParams{
		Note:           targetSpreadsheetName,
		Name:           targetSpreadsheetName,
		SpreadsheetUrl: "https://docs.google.com/spreadsheets/d/" + duplicateSpreadsheetResult.SpreadsheetId,
		SpreadsheetId:  duplicateSpreadsheetResult.SpreadsheetId,
		Password:       "",
		SheetName:      targetSheetName,
		SyncStrategy:   form.SyncStrategy,
	})

	if err != nil {
		return err
	}
	err = receiver.FormQuestionRepository.DeleteByFormID(newForm.ID)
	if err != nil {
		return err
	}

	err = receiver.saveQuestions(rawQuestions)
	if err != nil {
		return err
	}

	var formQuestionItems = make([]request.CreateFormQuestionItem, 0)
	for index, question := range rawQuestions {
		var answerRequired = false
		for _, formQuestion := range formQuestions {
			questionIndex := strings.Replace(formQuestion.QuestionId, fmt.Sprintf("%s_%s_", strings.ToUpper(form.Note), form.SpreadsheetId), "", -1)
			questionID := fmt.Sprintf("%s_%s_%s", strings.ToUpper(targetSpreadsheetName), duplicateSpreadsheetResult.SpreadsheetId, questionIndex)

			if questionID == question.QuestionId {
				answerRequired = formQuestion.AnswerRequired
			}
		}

		formQuestionItems = append(formQuestionItems, request.CreateFormQuestionItem{
			QuestionId:     question.QuestionId,
			Order:          index,
			AnswerRequired: answerRequired,
		})
	}
	_, err = receiver.FormQuestionRepository.CreateFormQuestions(newForm.ID, formQuestionItems)
	if err != nil {
		return err
	}

	//Duplicate submissions
	err = receiver.duplicateSubmissions(form, newForm, rq)
	if err != nil {
		return err
	}

	//Log the form url into device uploader spreadsheet
	defer func() {
		receiver.updateMemorySignUpFormIntoDeviceUploader(rq.DeviceId, duplicateSpreadsheetResult.SpreadsheetId)
	}()

	defer func() {
		receiver.cleanUpMemorySignUpForm(duplicateSpreadsheetResult.SpreadsheetId, "Sheet1")
	}()

	return nil
}

func (receiver *SubmitFormUseCase) saveQuestions(rawQuestions []entity.SQuestion) error {
	_, err := receiver.QuestionRepository.SaveQuestions(rawQuestions)

	return err
}

func (receiver *SubmitFormUseCase) updateMemorySignUpFormIntoDeviceUploader(deviceID string, spreadsheetID string) {
	setting, err := receiver.SettingRepository.GetSyncDevicesSettings()
	if err != nil {
		log.Error("Device uploader setting not found: ", err.Error())
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

	rows, err := receiver.Reader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Devices!K11:L",
	})
	if err != nil {
		log.Error("failed to read from sync devices sheet to find empty L cell")
		return
	}

	var rowNo int
	for rowNumber, row := range rows {
		if len(row) > 1 && row[1].(string) == deviceID {
			rowNo = rowNumber + 11
			break
		}
	}

	if rowNo == 0 {
		log.Error("No existing device found in sync devices sheet")
		monitor.SendMessageViaTelegram("Device ID: " + deviceID + " is now signing up but not found in sync devices uploader sheet")
		return
	}

	deviceData := make([][]interface{}, 0)
	deviceData = append(deviceData, []interface{}{"https://docs.google.com/spreadsheets/d/" + spreadsheetID})

	_, err = receiver.Writer.UpdateRange(sheet.WriteRangeParams{
		Range:     "Devices!I" + strconv.Itoa(rowNo),
		Rows:      deviceData,
		Dimension: "COLUMNS",
	}, spreadsheetId)

	if err != nil {
		log.Error("failed to update sync devices uploader")
		return
	}
}

func (receiver *SubmitFormUseCase) cleanUpMemorySignUpForm(spreadsheetId string, targetSheetName string) {
	_ = receiver.Writer.DeleteSheet(sheet.DeleteSheetParams{
		SpreadsheetID: spreadsheetId,
		SheetTitle:    targetSheetName,
	})
}

func (receiver *SubmitFormUseCase) UpdateSignUpMemoryForm(rq request.SubmitFormRequest) error {
	form, err := receiver.FormRepository.GetFormByQRCode(rq.QRCode)
	if err != nil {
		return err
	}

	submissionItems := make([]repository.SubmissionDataItem, 0)
	questions, err := receiver.QuestionRepository.GetQuestionsByIDs(Map(rq.Answers, func(answer request.Answer) string { return answer.QuestionId }))
	if err != nil {
		return fmt.Errorf("system cannot find questions for this form: %s", form.Name)
	}

	for _, answer := range rq.Answers {
		for _, question := range questions {
			if answer.QuestionId == question.QuestionId {
				var msg *repository.Messaging = nil
				if answer.Messaging != nil {
					msg = &repository.Messaging{
						Email:        answer.Messaging.Email,
						Value3:       answer.Messaging.Value3,
						MessageBox:   answer.Messaging.MessageBox,
						QuestionType: answer.Messaging.QuestionType,
					}
				}
				submissionItems = append(submissionItems, repository.SubmissionDataItem{
					QuestionId: question.QuestionId,
					Question:   question.Question,
					Answer:     answer.Answer,
					Messaging:  msg,
				})
			}
		}
	}

	submissionData := repository.SubmissionData{
		Items: submissionItems,
	}

	createSubmissionParams := repository.CreateSubmissionParams{
		FormId:         form.ID,
		DeviceId:       "SignedUp_At_" + time.Now().Format("20060102150405"),
		SubmissionData: submissionData,
		OpenedAt:       rq.OpenedAt,
	}
	err = receiver.SubmissionRepository.CreateSubmission(createSubmissionParams)
	if err != nil {
		log.Error("SubmitFormUseCase.UpdateSignUpMemoryForm", err)
		return errors.New("system cannot handle the submission")
	}

	createSubmissionParams = repository.CreateSubmissionParams{
		FormId:         form.ID,
		DeviceId:       "SignedUp_At_" + time.Now().Format("20060102150405"),
		SubmissionData: submissionData,
		OpenedAt:       rq.OpenedAt,
	}
	err = receiver.SubmissionRepository.DublicateSubmissions(createSubmissionParams)
	if err != nil {
		log.Error("SubmitFormUseCase.SubmitSignUpForm", err)
		return errors.New("system cannot handle the submission")
	}

	return nil
}

func (receiver *SubmitFormUseCase) duplicateSubmissions(systemSignUpForm entity.SForm, form *entity.SForm, rq request.SubmitFormRequest) error {

	submissionItems := make([]repository.SubmissionDataItem, 0)
	questions, err := receiver.QuestionRepository.
		GetQuestionsByIDs(Map(rq.Answers, func(answer request.Answer) string {
			return answer.QuestionId
		}),
		)
	if err != nil {
		return fmt.Errorf("system cannot find questions for this form: %s", form.Name)
	}

	for _, answer := range rq.Answers {
		for _, question := range questions {
			if answer.QuestionId == question.QuestionId {
				var msg *repository.Messaging = nil
				if answer.Messaging != nil {
					msg = &repository.Messaging{
						Email:        answer.Messaging.Email,
						Value3:       answer.Messaging.Value3,
						MessageBox:   answer.Messaging.MessageBox,
						QuestionType: answer.Messaging.QuestionType,
					}
				}

				questionIndex := strings.Replace(answer.QuestionId, fmt.Sprintf("%s_%s_", strings.ToUpper(systemSignUpForm.Note), systemSignUpForm.SpreadsheetId), "", -1)
				questionID := fmt.Sprintf("%s_%s_%s", strings.ToUpper(form.Note), form.SpreadsheetId, questionIndex)

				submissionItems = append(submissionItems, repository.SubmissionDataItem{
					QuestionId: questionID,
					Question:   question.Question,
					Answer:     answer.Answer,
					Messaging:  msg,
				})
			}
		}
	}

	submissionData := repository.SubmissionData{
		Items: submissionItems,
	}

	// setting, err := receiver.SettingRepository.GetRegistrationSubmissionSetting()
	// if err != nil {
	// 	return err
	// }

	// if setting == nil {
	// 	return errors.New("registration submission setting is not set")
	// }

	// type summarySetting struct {
	// 	SpreadsheetId string `json:"spreadsheet_id"`
	// }

	// var summary summarySetting
	// err = json.Unmarshal(setting.Settings, &summary)
	// if err != nil {
	// 	return err
	// }

	// re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	// match := re.FindStringSubmatch(summary.SpreadsheetId)

	// if len(match) < 2 {
	// 	return errors.New("invalid spreadsheet url from sign up submission setting")
	// }

	createSubmissionParmas := repository.CreateSubmissionParams{
		FormId:         form.ID,
		DeviceId:       "SignedUp_At_" + time.Now().Format("20060102150405"),
		SubmissionData: submissionData,
		OpenedAt:       rq.OpenedAt,
	}
	err = receiver.SubmissionRepository.DublicateSubmissions(createSubmissionParmas)
	if err != nil {
		log.Error("SubmitFormUseCase.duplicateSubmissions", err)
		return errors.New("system cannot handle the submission")
	}

	return nil
}
