package usecase

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"regexp"
	"sen-global-api/config"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/parameters"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/job"
	"sen-global-api/pkg/monitor"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"
	"time"
)

type FormsUploaderIndex int

const (
	FormsUploaderIndexFirst FormsUploaderIndex = iota
	FormsUploaderIndexSecond
	FormsUploaderIndexThird
	FormsUploaderIndexFourth
	FormsUploaderIndexFifth //Sign Up Forms
)

type ImportFormsUseCase struct {
	FormRepository                  *repository.FormRepository
	QuestionRepository              *repository.QuestionRepository
	FormQuestionRepository          *repository.FormQuestionRepository
	SpreadsheetReader               *sheet.Reader
	SpreadsheetWriter               *sheet.Writer
	SettingRepository               *repository.SettingRepository
	DefaultCronJobIntervalInMinutes uint8
	TimeMachine                     *job.TimeMachine
	config.AppConfig
}

func (receiver *ImportFormsUseCase) SyncForms(req request.ImportFormRequest) error {
	monitor.SendMessageViaTelegram(fmt.Sprintf("[INFO][SYNC] Start sync Forms %s with interval %d", req.SpreadsheetUrl, req.Interval))
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}
	if req.Interval == 0 {
		return nil
	}

	spreadsheetId := match[1]
	monitor.LogGoogleAPIRequestImportForm()

	sheets, err := receiver.SpreadsheetReader.GetSheets(spreadsheetId)

	if err != nil {
		log.Error(err)
		return err
	}

	if len(sheets) == 0 {
		return fmt.Errorf("no sheet found")
	}

	for _, sheetName := range sheets {
		if strings.HasPrefix(strings.ToLower(sheetName), "[up]") == false {
			continue
		}
		values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
			SpreadsheetId: spreadsheetId,
			ReadRange:     sheetName + `!` + receiver.AppConfig.Google.FirstColumn + strconv.Itoa(receiver.Google.FirstRow+2) + `:AC`,
		})
		if err != nil {
			log.Error(err)
			return err
		}
		for rowNo, row := range values {
			if len(row) >= 4 && cap(row) >= 4 {
				formStatus, err := value.GetImportSpreadsheetStatusFromString(row[3].(string))
				if err != nil {
					return err
				}
				switch formStatus {
				case value.ImportSpreadsheetStatusDeleted:
					err = receiver.FormRepository.DeleteFormByNote(row[0].(string))
					if err != nil {
						log.Error(err)
					} else {
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"DELETED", time.Now().Format("2006-01-02 15:04:05"), ""}},
						}, spreadsheetId)
						if err != nil {
							log.Debug("Row No: ", rowNo)
							log.Error(err)
						}
					}
				case value.ImportSpreadsheetStatusDeactivate:
					err = receiver.FormRepository.DeleteFormByNote(row[0].(string))
					if err != nil {
						log.Error(err)
					} else {
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"DEACTIVATED", time.Now().Format("2006-01-02 15:04:05"), ""}},
						}, spreadsheetId)
						if err != nil {
							log.Debug("Row No: ", rowNo)
							log.Error(err)
						}
					}
				case value.ImportSpreadsheetStatusPending:
					continue
				case value.ImportSpreadsheetStatusSkip:
					continue
				case value.ImportSpreadsheetStatusNew:
					if row[0].(string) == "" || row[1].(string) == "" {
						//Skip row if code or spreadsheet url is empty
						continue
					}
					var submissionType value.SubmissionType = value.SubmissionTypeValues
					var submissionSheetId string = ""
					if len(row) >= 15 {
						submissionType = value.GetSubmissionTypeFromString(row[14].(string))
						if len(row) >= 16 {
							//Just required column Z (#15) in case of submission type is in (2,3,5)
							switch submissionType {
							case value.SubmissionTypeBoth, value.SubmissionTypeQrCode, value.SubmissionTypeTeacherAndQRCode:
								re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
								match := re.FindStringSubmatch(row[15].(string))

								if len(match) < 2 {
									continue
								} else {
									submissionSheetId = match[1]
								}
							default:
								break
							}
						}
					}
					var tabName string = ""
					if len(row) >= 17 {
						tabName = row[16].(string)
					}
					outputSheetName := "Answers"
					if len(row) >= 18 {
						outputSheetName = row[17].(string)
					}
					syncStrategy := value.FormSyncStrategyOnSubmit
					if len(row) >= 19 {
						syncStrategy = value.GetFormSyncStrategyFromString(row[18].(string))
					}
					importErr, reason := receiver.importForm(row[0].(string), row[1].(string), row[2].(string), row[3].(string), submissionType, submissionSheetId, tabName, outputSheetName, syncStrategy)
					if importErr != nil {
						log.Error(importErr)
						monitor.LogGoogleAPIRequestImportForm()
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
						}, spreadsheetId)
						if err != nil {
							log.Debug("Row No: ", rowNo)
							log.Error(err)
						}
					} else {
						monitor.LogGoogleAPIRequestImportForm()
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
						}, spreadsheetId)
						if err != nil {
							log.Debug("Row No: ", rowNo)
							log.Error(err)
						}
					}
				}
			}
		}
	}

	return nil
}

func (receiver *ImportFormsUseCase) ImportForms(req request.ImportFormRequest, uploaderIndex FormsUploaderIndex) error {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}
	if req.Interval == 0 {
		req.AutoImport = false
	}

	switch uploaderIndex {
	case FormsUploaderIndexFirst:
		err := receiver.SettingRepository.UpdateFormSetting(req)
		if err != nil {
			return err
		}
	case FormsUploaderIndexSecond:
		err := receiver.SettingRepository.UpdateFormSetting2(req)
		if err != nil {
			return err
		}
	case FormsUploaderIndexThird:
		err := receiver.SettingRepository.UpdateFormSetting3(req)
		if err != nil {
			return err
		}
	case FormsUploaderIndexFourth:
		err := receiver.SettingRepository.UpdateFormSetting4(req)
		if err != nil {
			return err
		}
	case FormsUploaderIndexFifth:
		err := receiver.SettingRepository.UpdateSignUpFormSetting(req)
		if err != nil {
			return err
		}

		return receiver.importSignUpForms(req)
	}

	spreadsheetId := match[1]
	monitor.LogGoogleAPIRequestImportForm()

	sheets, err := receiver.SpreadsheetReader.GetSheets(spreadsheetId)

	if err != nil {
		log.Error(err)
		return err
	}

	if len(sheets) == 0 {
		return fmt.Errorf("no sheet found")
	}

	for _, sheetName := range sheets {
		if strings.HasPrefix(strings.ToLower(sheetName), "[up]") == false {
			continue
		}
		values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
			SpreadsheetId: spreadsheetId,
			ReadRange:     sheetName + `!` + receiver.AppConfig.Google.FirstColumn + strconv.Itoa(receiver.Google.FirstRow+2) + `:AC`,
		})
		if err != nil {
			log.Error(err)
			return err
		}
		for rowNo, row := range values {
			if len(row) >= 4 && cap(row) >= 4 {
				formStatus, err := value.GetImportSpreadsheetStatusFromString(row[3].(string))
				if err != nil {
					return err
				}
				switch formStatus {
				case value.ImportSpreadsheetStatusDeleted:
					err = receiver.FormRepository.DeleteFormByNote(row[0].(string))
					if err != nil {
						log.Error(err)
					} else {
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"DELETED", time.Now().Format("2006-01-02 15:04:05"), ""}},
						}, spreadsheetId)
						if err != nil {
							log.Debug("Row No: ", rowNo)
							log.Error(err)
						}
					}
				case value.ImportSpreadsheetStatusDeactivate:
					err = receiver.FormRepository.DeleteFormByNote(row[0].(string))
					if err != nil {
						log.Error(err)
					} else {
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"DEACTIVATED", time.Now().Format("2006-01-02 15:04:05"), ""}},
						}, spreadsheetId)
						if err != nil {
							log.Debug("Row No: ", rowNo)
							log.Error(err)
						}
					}
				case value.ImportSpreadsheetStatusPending:
					continue
				case value.ImportSpreadsheetStatusSkip:
					continue
				case value.ImportSpreadsheetStatusNew:
					if row[0].(string) == "" || row[1].(string) == "" {
						//Skip row if code or spreadsheet url is empty
						continue
					}
					var submissionType value.SubmissionType = value.SubmissionTypeValues
					var submissionSheetId string = ""
					if len(row) >= 15 {
						submissionType = value.GetSubmissionTypeFromString(row[14].(string))
						if len(row) >= 16 {
							//Just required column Z (#15) in case of submission type is in (2,3,5)
							switch submissionType {
							case value.SubmissionTypeBoth, value.SubmissionTypeQrCode, value.SubmissionTypeTeacherAndQRCode:
								re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
								match := re.FindStringSubmatch(row[15].(string))

								if len(match) < 2 {
									continue
								} else {
									submissionSheetId = match[1]
								}
							default:
								break
							}
						}
					}
					var tabName string = ""
					if len(row) >= 17 {
						tabName = row[16].(string)
					}
					outputSheetName := "Answers"
					if len(row) >= 18 {
						outputSheetName = row[17].(string)
					}
					syncStrategy := value.FormSyncStrategyOnSubmit
					if len(row) >= 19 {
						syncStrategy = value.GetFormSyncStrategyFromString(row[18].(string))
					}
					importErr, reason := receiver.importForm(row[0].(string), row[1].(string), row[2].(string), row[3].(string), submissionType, submissionSheetId, tabName, outputSheetName, syncStrategy)
					if importErr != nil {
						log.Error(importErr)
						monitor.LogGoogleAPIRequestImportForm()
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
						}, spreadsheetId)
						if err != nil {
							log.Debug("Row No: ", rowNo)
							log.Error(err)
						}
					} else {
						monitor.LogGoogleAPIRequestImportForm()
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
						}, spreadsheetId)
						if err != nil {
							log.Debug("Row No: ", rowNo)
							log.Error(err)
						}
					}
				}
			}
		}
	}

	var interval uint64 = 0
	if req.AutoImport != false {
		interval = req.Interval
	}

	switch uploaderIndex {
	case FormsUploaderIndexFirst:
		receiver.TimeMachine.ScheduleSyncForms(interval)
	case FormsUploaderIndexSecond:
		receiver.TimeMachine.ScheduleSyncForms2(interval)
	case FormsUploaderIndexThird:
		receiver.TimeMachine.ScheduleSyncForms3(interval)
	case FormsUploaderIndexFourth:
		receiver.TimeMachine.ScheduleSyncForms4(interval)
	case FormsUploaderIndexFifth:
		log.Error("FormsUploaderIndexFifth must not sync here")
	}

	return nil
}

func (receiver *ImportFormsUseCase) importSignUpForms(req request.ImportFormRequest) error {
	//Get spreadsheet data at sheet Forms

	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)
	signUpFormsSpreadsheetId := match[1]

	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: signUpFormsSpreadsheetId,
		ReadRange:     "Forms" + `!K11:M`,
	})
	if err != nil || values == nil {
		log.Error(err)
		return errors.New("Error reading sign up forms spreadsheet: " + req.SpreadsheetUrl + " - " + err.Error())
	}

	for _, row := range values {
		if len(row) < 3 {
			continue
		}

		code := row[0].(string)
		sheetName := row[2].(string)

		if sheetName == "" || code == "" {
			continue
		}

		url := row[1].(string)
		//validate spreadsheet url
		if !strings.Contains(url, "https://docs.google.com/spreadsheets/d/") {
			continue
		}

		_, err = receiver.importSignUpForm(url, code, sheetName)
		if err != nil {
			monitor.SendMessageViaTelegram(
				"[ERROR]: Error importing sign up forms: "+err.Error(),
				"Detail: form code: "+code,
				"Sheet Name: "+sheetName,
				"Spreadsheet Url: "+url,
			)
		}
	}

	return nil
}

func (c *ImportFormsUseCase) importSignUpForm(spreadsheetUrl, note, sheetNameToRead string) (entity.SForm, error) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(spreadsheetUrl)

	if len(match) < 2 {
		log.Error("Import Sign Up Form Invalid spreadsheet url: ", spreadsheetUrl)
		return entity.SForm{},
			fmt.Errorf(fmt.Sprintf("Import Sign Up Form Invalid spreadsheet url: %s", spreadsheetUrl))
	}

	spreadsheetId := match[1]
	monitor.LogGoogleAPIRequestImportForm()
	values, err := c.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     sheetNameToRead + `!J11` + `:Q`,
	})
	if err != nil || values == nil {
		log.Error(fmt.Sprintf("Error reading spreadsheet: %s - note : %s", err.Error(), note))
		return entity.SForm{}, err
	}

	var rawQuestions = make([]parameters.RawQuestion, 0)
	var formName = ""
	for index, row := range values {
		if index == 0 && len(row) > 1 && cap(row) > 1 {
			formName = row[1].(string)
			continue
		} else if len(row) > 4 && cap(row) > 4 && index > 1 && row[2].(string) != "" {
			additionalInfo := ""
			if len(row) > 6 {
				additionalInfo = row[6].(string)
			}
			required := "false"
			if len(row) > 5 {
				required = row[5].(string)
			}
			var uniqueId *string = nil
			if row[0].(string) != "" {
				id := row[0].(string)
				uniqueId = &id
			}
			item := parameters.RawQuestion{
				QuestionId:        strings.ToUpper(note) + "_" + spreadsheetId + "_" + strconv.Itoa(index-1),
				Question:          row[3].(string),
				Type:              row[2].(string),
				Attributes:        strings.ReplaceAll(row[4].(string), "\n", ""),
				AnswerRequired:    required,
				AdditionalOptions: additionalInfo,
				Status:            "1",
				RowNumber:         index + 1,
				QuestionUniqueId:  uniqueId,
			}
			rawQuestions = append(rawQuestions, item)
		}
	}

	f, err, msg := c.CreateSignUpForm(parameters.SaveFormParams{
		Note:              note,
		Name:              formName,
		SpreadsheetUrl:    spreadsheetUrl,
		SpreadsheetId:     spreadsheetId,
		Password:          "",
		RawQuestions:      rawQuestions,
		SubmissionType:    value.SubmissionTypeSignUpRegistration,
		SubmissionSheetId: "",
		SheetName:         sheetNameToRead,
		OutputSheetName:   "",
		SyncStrategy:      value.FormSyncStrategyOnSubmit,
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error creating form: %s - note : %s", err.Error(), note))
		return entity.SForm{}, err
	}

	log.Warning(msg)

	return *f, nil
}

func (receiver *ImportFormsUseCase) importForm(code string, spreadsheetUrl string, password string, status string, submissionType value.SubmissionType, submissionSheetId string, sheetName string, outputSheetName string, syncStrategy value.FormSyncStrategy) (error, string) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(spreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url"), "Invalid spreadsheet url"
	}

	spreadsheetId := match[1]
	var sheetNameToRead string = "Questions"
	if sheetName != "" {
		sheetNameToRead = sheetName
	}
	monitor.LogGoogleAPIRequestImportForm()
	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     sheetNameToRead + `!I11:Q`,
	})
	if err != nil || values == nil {
		log.Error(err)
		return err, fmt.Sprintf("Cannot read tab %s from %s ERROR %s", sheetNameToRead, spreadsheetUrl, err.Error())
	}

	var rawQuestions = make([]parameters.RawQuestion, 0)
	var formName = ""
	for index, row := range values {
		if index == 0 && len(row) > 2 && cap(row) > 2 {
			formName = row[2].(string)
			continue
		} else if len(row) > 5 && cap(row) > 5 && index > 1 && row[3].(string) != "" {
			additionalInfo := ""
			if len(row) > 7 {
				additionalInfo = row[7].(string)
			}
			required := "false"
			if len(row) > 6 {
				required = row[6].(string)
			}
			var enabled value.QuestionForMobile = value.QuestionForMobile_Enabled
			if strings.ToUpper(row[0].(string)) == "LOCK" {
				enabled = value.QuestionForMobile_Disabled
			}
			item := parameters.RawQuestion{
				QuestionId:        strings.ToUpper(code) + "_" + spreadsheetId + "_" + row[2].(string),
				Question:          row[4].(string),
				Type:              row[3].(string),
				Attributes:        strings.ReplaceAll(row[5].(string), "\n", ""),
				AnswerRequired:    required,
				AdditionalOptions: additionalInfo,
				Status:            "1",
				RowNumber:         index + 1,
				EnableOnMobile:    enabled,
			}
			rawQuestions = append(rawQuestions, item)
		}
	}

	_, err, reason := receiver.saveForm(parameters.SaveFormParams{
		Note:              code,
		Name:              formName,
		SpreadsheetUrl:    spreadsheetUrl,
		SpreadsheetId:     spreadsheetId,
		Password:          password,
		RawQuestions:      rawQuestions,
		SubmissionType:    submissionType,
		SubmissionSheetId: submissionSheetId,
		SheetName:         sheetNameToRead,
		OutputSheetName:   outputSheetName,
		SyncStrategy:      syncStrategy,
	})

	return err, reason
}

func (receiver *ImportFormsUseCase) CreateSignUpForm(params parameters.SaveFormParams) (*entity.SForm, error, string) {
	return receiver.saveForm(params)
}

func (receiver *ImportFormsUseCase) saveForm(params parameters.SaveFormParams) (*entity.SForm, error, string) {
	err := receiver.QuestionRepository.DeleteQuestionsFormNote(params.Note)
	if err != nil {
		return nil, err, "System Error: cannot delete questions belong to this form: " + params.Note
	}

	questions, err, invalidQuestions := receiver.saveQuestions(params.RawQuestions)
	var reason string
	if len(invalidQuestions) > 0 {
		reason = "Invalid questions: "
		for _, question := range invalidQuestions {
			reason += "Row No:" + strconv.Itoa(question.RowNumber) + ": " + question.Reason + ", "
		}
	}
	if err != nil {
		return nil, err, "" + err.Error() + " " + reason
	}
	log.Debug(questions)

	form, err := receiver.createForm(questions, params)

	return form, err, reason
}

type InvalidQuestionRow struct {
	RowNumber int
	Reason    string
}

func (receiver *ImportFormsUseCase) saveQuestions(rawQuestions []parameters.RawQuestion) ([]entity.SQuestion, error, []InvalidQuestionRow) {
	var params = make([]repository.CreateQuestionParams, 0)
	var invalidQuestions = make([]InvalidQuestionRow, 0)
	for i, rawQuestion := range rawQuestions {
		questionType, err := value.GetQuestionType(rawQuestion.Type)
		if err != nil {
			log.Info(fmt.Sprintf("Invalid question type: %s - %v", rawQuestion.Type, rawQuestion))
			invalidQuestions = append(invalidQuestions, InvalidQuestionRow{
				RowNumber: i + 2,
				Reason:    fmt.Sprintf("Invalid question type: %s - %v", rawQuestion.Type, rawQuestion),
			})
			continue
		}

		status, err := GetStatusFromString(rawQuestion.Status)
		if err != nil {
			invalidQuestions = append(invalidQuestions, InvalidQuestionRow{
				RowNumber: i + 2,
				Reason:    "Invalid question status: " + rawQuestion.Status,
			})
			continue
		}

		attString, err := UnmarshalAttributes(rawQuestion, questionType)
		if err != nil {
			invalidQuestions = append(invalidQuestions, InvalidQuestionRow{
				RowNumber: i + 2,
				Reason:    "Invalid question attribute: " + rawQuestion.Attributes,
			})
			continue
		}

		param := repository.CreateQuestionParams{
			QuestionId:     rawQuestion.QuestionId,
			QuestionName:   rawQuestion.Question,
			QuestionType:   strings.ToLower(rawQuestion.Type),
			Question:       rawQuestion.Question,
			Attributes:     attString,
			Status:         value.GetRawStatusValue(status),
			Set:            rawQuestion.Attributes,
			EnableOnMobile: rawQuestion.EnableOnMobile,
			QuestionUniqueId: rawQuestion.QuestionUniqueId,
		}
		params = append(params, param)
	}

	if len(params) == 0 {
		return nil, errors.New("this form does not have any valid format questions"), make([]InvalidQuestionRow, 0)
	}

	questions, err := receiver.QuestionRepository.Create(params)

	return questions, err, invalidQuestions
}

func GetStatusFromString(status string) (value.Status, error) {
	switch strings.ToLower(status) {
	case "true":
		return value.Active, nil
	case "false":
		return value.Inactive, nil
	default:
		return value.Active, nil
	}
}

func UnmarshalAttributes(rawQuestion parameters.RawQuestion, questionType value.QuestionType) (string, error) {

	switch questionType {
	case value.QuestionTime:
		return "{}", nil
	case value.QuestionDate:
		return "{}", nil
	case value.QuestionDateTime:
		return "{}", nil
	case value.QuestionDurationForward:
		return "{}", nil
	case value.QuestionDurationBackward:
		//TODO: validate attributes
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionScale:
		rawValues := strings.Split(rawQuestion.Attributes, ";")
		if len(rawValues) < 2 {
			return "", errors.New("scale question data is invalid " + rawQuestion.Attributes)
		}
		totalValuesInString := strings.Split(rawValues[0], ":")
		if len(totalValuesInString) < 2 {
			return "", errors.New("scale question data is invalid " + rawQuestion.Attributes)
		}
		stepValueInString := strings.Split(rawValues[1], ":")
		if len(stepValueInString) < 2 {
			return "", errors.New("scale question data is invalid " + rawQuestion.Attributes)
		}
		totalValues, err := strconv.Atoi(totalValuesInString[1])
		if err != nil {
			return "", errors.New("scale question data is invalid " + err.Error())
		}
		stepValue, err := strconv.Atoi(stepValueInString[1])
		if err != nil {
			return "", errors.New("scale question data is invalid " + err.Error())
		}
		return "{\"number\" : " + strconv.Itoa(totalValues) + ", \"steps\": " + strconv.Itoa(stepValue) + "}", nil
	case value.QuestionQRCode:
		return "{}", nil
	case value.QuestionSelection:
		rawOptions := strings.Split(rawQuestion.Attributes, ";")
		//`{"options": [{"name": "red"}, { "name": "green"}, {"name" : "blue"}]}`,
		type Option struct {
			Name string `json:"name"`
		}
		type Options struct {
			Options []Option `json:"options"`
		}
		var optionsList = make([]Option, 0)
		for _, op := range rawOptions {
			if op == "" {
				return "", errors.New("invalid options")
			}
			optionsList = append(optionsList, Option{Name: op})
		}
		options := Options{Options: optionsList}
		result, err := json.Marshal(options)
		if err != nil {
			return "", err
		}
		return string(result), nil
	case value.QuestionText:
		return "{}", nil
	case value.QuestionCount:
		return "{}", nil
	case value.QuestionNumber:
		return "{}", nil
	case value.QuestionPhoto:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionButtonCount:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionMultipleChoice:
		rawOptions := strings.Split(rawQuestion.Attributes, ";")
		type Option struct {
			Name string `json:"name"`
		}
		type Options struct {
			Options []Option `json:"options"`
		}
		var optionsList = make([]Option, 0)
		for _, op := range rawOptions {
			if op == "" {
				return "", errors.New("invalid options")
			}
			optionsList = append(optionsList, Option{Name: op})
		}
		options := Options{Options: optionsList}
		result, err := json.Marshal(options)
		if err != nil {
			return "", err
		}
		return string(result), nil
	case value.QuestionSingleChoice:
		rawOptions := strings.Split(rawQuestion.Attributes, ";")
		type Option struct {
			Name string `json:"name"`
		}
		type Options struct {
			Options []Option `json:"options"`
		}
		var optionsList = make([]Option, 0)
		for _, op := range rawOptions {
			if op == "" {
				return "", errors.New("invalid options")
			}
			optionsList = append(optionsList, Option{Name: op})
		}
		options := Options{Options: optionsList}
		result, err := json.Marshal(options)
		if err != nil {
			return "", err
		}
		return string(result), nil

	case value.QuestionButtonList:
		re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
		match := re.FindStringSubmatch(rawQuestion.Attributes)

		if len(match) < 2 {
			return "", errors.New("invalid google sheet url")
		}

		return `{"spreadsheet_id" : "` + match[1] + `"}`, nil

	case value.QuestionMessageBox:
		message := strings.Replace(rawQuestion.Attributes, "\n", "\\n", -1)
		jsonMsg := `{"value": "` + message + `"}`
		return jsonMsg, nil
	case value.QuestionShowPic:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionButton:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionPlayVideo:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionQRCodeFront:
		return "{}", nil
	case value.QuestionRandomizer:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionDocument:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionQRCodeGenerator:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionCodeCounting:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionChoiceToggle:
		rawOptions := strings.Split(rawQuestion.Attributes, ";")
		type Option struct {
			Name string `json:"name"`
		}
		type Options struct {
			Options []Option `json:"options"`
		}
		var optionsList = make([]Option, 0)
		for _, op := range rawOptions {
			if op == "" {
				return "", errors.New("invalid options")
			}
			optionsList = append(optionsList, Option{Name: op})
		}
		options := Options{Options: optionsList}
		result, err := json.Marshal(options)
		if err != nil {
			return "", err
		}
		return string(result), nil

	case value.QuestionSection:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil

	case value.QuestionFormSection:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil

	case value.QuestionFormSendImmediately:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionSignature:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionWeb:
		attInbase64 := base64.StdEncoding.EncodeToString([]byte(rawQuestion.Attributes))
		return `{"value": "` + attInbase64 + `"}`, nil
	case value.QuestionSignUpPreSetValue1:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionSignUpPreSetValue2:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionSignUpPreSetValue3:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionDraggableList:
		rawOptions := strings.Split(rawQuestion.Attributes, ";")
		type Option struct {
			Name string `json:"name"`
		}
		type Options struct {
			Options []Option `json:"options"`
		}
		var optionsList = make([]Option, 0)
		for _, op := range rawOptions {
			if op == "" {
				return "", errors.New("invalid options")
			}
			optionsList = append(optionsList, Option{Name: op})
		}
		options := Options{Options: optionsList}
		result, err := json.Marshal(options)
		if err != nil {
			return "", err
		}
		return string(result), nil
	case value.QuestionSendMessage:
		type Msg struct {
			Email          []string `json:"email"`
			Value3         []string `json:"value3"`
			ShowMessageBox bool     `json:"showMessageBox"`
		}
		type Messaging struct {
			Data Msg `json:"messaging"`
		}

		var msg Msg
		err := json.Unmarshal([]byte(rawQuestion.Attributes), &msg)
		if err != nil {
			return "", err
		}
		messaging := Messaging{Data: msg}
		result, err := json.Marshal(messaging)
		if err != nil {
			return "", err
		}
		return string(result), nil
	case value.QuestionSendNotification:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	default:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	}

}

func (receiver *ImportFormsUseCase) createForm(questions []entity.SQuestion, params parameters.SaveFormParams) (*entity.SForm, error) {
	form, err := receiver.FormRepository.SaveForm(params)
	if err != nil {
		return nil, err
	}
	var formQuestions = make([]request.CreateFormQuestionItem, 0)
	for _, question := range questions {
		var order = 0
		var answerRequired = false
		for _, rq := range params.RawQuestions {
			if rq.QuestionId == question.QuestionId {
				order = rq.RowNumber
				answerRequired = strings.ToLower(rq.AnswerRequired) == "true"
			}
		}

		formQuestions = append(formQuestions, request.CreateFormQuestionItem{
			QuestionId:     question.QuestionId,
			Order:          order,
			AnswerRequired: answerRequired,
		})
	}
	_, err = receiver.FormQuestionRepository.CreateFormQuestions(form.FormId, formQuestions)
	if err != nil {
		return nil, err
	}
	return form, nil
}

// ImportFormsPartially Fetch form data from google sheet at sheet name tabName and save to database
func (receiver *ImportFormsUseCase) ImportFormsPartially(url string, sheetName string) error {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(url)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}
	spreadsheetId := match[1]
	monitor.LogGoogleAPIRequestImportForm()

	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     sheetName + `!` + receiver.AppConfig.Google.FirstColumn + strconv.Itoa(receiver.Google.FirstRow+2) + `:AC`,
	})
	if err != nil {
		log.Error(err)
		return errors.New(fmt.Sprintf("Error reading spreadsheet: %s", err.Error()))
	}
	for rowNo, row := range values {
		if len(row) >= 4 && cap(row) >= 4 {
			formStatus, err := value.GetImportSpreadsheetStatusFromString(row[3].(string))
			if err != nil {
				return err
			}
			switch formStatus {
			case value.ImportSpreadsheetStatusDeleted:
				err = receiver.FormRepository.DeleteFormByNote(row[0].(string))
				if err != nil {
					log.Error(err)
				} else {
					_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
						Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
						Dimension: "ROWS",
						Rows:      [][]interface{}{{"DELETED", time.Now().Format("2006-01-02 15:04:05"), ""}},
					}, spreadsheetId)
					if err != nil {
						log.Debug("Row No: ", rowNo)
						log.Error(err)
					}
				}
			case value.ImportSpreadsheetStatusDeactivate:
				err = receiver.FormRepository.DeleteFormByNote(row[0].(string))
				if err != nil {
					log.Error(err)
				} else {
					_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
						Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
						Dimension: "ROWS",
						Rows:      [][]interface{}{{"DEACTIVATED", time.Now().Format("2006-01-02 15:04:05"), ""}},
					}, spreadsheetId)
					if err != nil {
						log.Debug("Row No: ", rowNo)
						log.Error(err)
					}
				}
			case value.ImportSpreadsheetStatusPending:
				continue
			case value.ImportSpreadsheetStatusSkip:
				continue
			case value.ImportSpreadsheetStatusNew:
				if row[0].(string) == "" || row[1].(string) == "" {
					//Skip row if code or spreadsheet url is empty
					continue
				}
				var submissionType value.SubmissionType = value.SubmissionTypeValues
				var submissionSheetId string = ""
				if len(row) >= 15 {
					submissionType = value.GetSubmissionTypeFromString(row[14].(string))
					if len(row) >= 16 {
						//Just required column Z (#15) in case of submission type is in (2,3,5)
						switch submissionType {
						case value.SubmissionTypeBoth, value.SubmissionTypeQrCode, value.SubmissionTypeTeacherAndQRCode: //5
							re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
							match := re.FindStringSubmatch(row[15].(string))

							if len(match) < 2 {
								continue
							} else {
								submissionSheetId = match[1]
							}
						default:
							break
						}
					}
				}
				var tabName string = ""
				if len(row) >= 17 {
					tabName = row[16].(string)
				}
				outputSheetName := "Answers"
				if len(row) >= 18 {
					outputSheetName = row[17].(string)
				}
				syncStrategy := value.FormSyncStrategyOnSubmit
				if len(row) >= 19 {
					syncStrategy = value.GetFormSyncStrategyFromString(row[18].(string))
				}
				importErr, reason := receiver.importForm(row[0].(string), row[1].(string), row[2].(string), row[3].(string), submissionType, submissionSheetId, tabName, outputSheetName, syncStrategy)
				if importErr != nil {
					log.Error(importErr)
					monitor.LogGoogleAPIRequestImportForm()
					_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
						Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
						Dimension: "ROWS",
						Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
					}, spreadsheetId)
					if err != nil {
						log.Debug("Row No: ", rowNo)
						log.Error(err)
					}
				} else {
					monitor.LogGoogleAPIRequestImportForm()
					_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
						Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
						Dimension: "ROWS",
						Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
					}, spreadsheetId)
					if err != nil {
						log.Debug("Row No: ", rowNo)
						log.Error(err)
					}
				}
			}
		}
	}

	return nil
}
