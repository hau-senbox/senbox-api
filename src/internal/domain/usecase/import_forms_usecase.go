package usecase

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
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

	"github.com/google/uuid"

	log "github.com/sirupsen/logrus"
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
	RoleOrgSignUpRepo               *repository.RoleOrgSignUpRepository
	DefaultCronJobIntervalInMinutes uint8
	TimeMachine                     *job.TimeMachine
	config.AppConfig
}

func (receiver *ImportFormsUseCase) SyncForms(req request.ImportFormRequest) error {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		return fmt.Errorf("invalid spreadsheet url")
	}
	if req.Interval == 0 {
		return nil
	}

	spreadsheetID := match[1]
	monitor.LogGoogleAPIRequestImportForm()

	sheets, err := receiver.SpreadsheetReader.GetSheets(spreadsheetID)

	if err != nil {
		log.Error(err)
		return err
	}

	if len(sheets) == 0 {
		return fmt.Errorf("no sheet found")
	}

	for _, sheetName := range sheets {
		if !strings.HasPrefix(strings.ToLower(sheetName), "[up]") {
			continue
		}
		values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
			SpreadsheetID: spreadsheetID,
			ReadRange:     sheetName + `!` + receiver.Google.FirstColumn + strconv.Itoa(receiver.Google.FirstRow+2) + `:AC`,
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
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"DELETED", time.Now().Format("2006-01-02 15:04:05"), ""}},
						}, spreadsheetID)
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
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"DEACTIVATED", time.Now().Format("2006-01-02 15:04:05"), ""}},
						}, spreadsheetID)
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
					if len(row) >= 15 {
						if len(row) >= 16 {
							//Just required column Z (#15) in case of submission type is in (2,3,5)
							re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
							match := re.FindStringSubmatch(row[15].(string))

							if len(match) < 2 {
								continue
							}
						}
					}
					tabName := ""
					if len(row) >= 17 {
						tabName = row[16].(string)
					}
					reason, importErr := receiver.importForm(row[0].(string), row[1].(string), row[2].(string), tabName)
					if importErr != nil {
						log.Error(importErr)
						monitor.LogGoogleAPIRequestImportForm()
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
						}, spreadsheetID)
						if err != nil {
							log.Debug("Row No: ", rowNo)
							log.Error(err)
						}
					} else {
						monitor.LogGoogleAPIRequestImportForm()
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
						}, spreadsheetID)
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

	spreadsheetID := match[1]
	monitor.LogGoogleAPIRequestImportForm()

	sheets, err := receiver.SpreadsheetReader.GetSheets(spreadsheetID)

	if err != nil {
		log.Error(err)
		return err
	}

	if len(sheets) == 0 {
		return fmt.Errorf("no sheet found")
	}

	for _, sheetName := range sheets {
		if !strings.HasPrefix(strings.ToLower(sheetName), "[up]") {
			continue
		}
		values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
			SpreadsheetID: spreadsheetID,
			ReadRange:     sheetName + `!` + receiver.Google.FirstColumn + strconv.Itoa(receiver.Google.FirstRow+2) + `:AC`,
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
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"DELETED", time.Now().Format("2006-01-02 15:04:05"), ""}},
						}, spreadsheetID)
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
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"DEACTIVATED", time.Now().Format("2006-01-02 15:04:05"), ""}},
						}, spreadsheetID)
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
					if len(row) >= 15 {
						if len(row) >= 16 {
							// //Just required column Z (#15) in case of submission type is in (2,3,5)
							// re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
							// match := re.FindStringSubmatch(row[15].(string))

							// if len(match) < 2 {
							// 	break
							// }
						}
					}
					tabName := ""
					if len(row) >= 17 {
						tabName = row[16].(string)
					}
					reason, importErr := receiver.importForm(row[0].(string), row[1].(string), row[2].(string), tabName)
					if importErr != nil {
						log.Error(importErr)
						monitor.LogGoogleAPIRequestImportForm()
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
						}, spreadsheetID)
						if err != nil {
							log.Debug("Row No: ", rowNo)
							log.Error(err)
						}
					} else {
						monitor.LogGoogleAPIRequestImportForm()
						_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
							Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.Google.FirstRow+2) + ":Q",
							Dimension: "ROWS",
							Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
						}, spreadsheetID)
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
	if req.AutoImport {
		interval = req.Interval
	}

	switch uploaderIndex {
	case FormsUploaderIndexFirst:
		receiver.TimeMachine.ScheduleSyncForms(interval)
	// case FormsUploaderIndexSecond:
	// 	receiver.TimeMachine.ScheduleSyncForms2(interval)
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
	signUpFormsSpreadsheetID := match[1]

	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetID: signUpFormsSpreadsheetID,
		ReadRange:     "Forms" + `!K11:P`,
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

		// luu vao bang RoleOrgSignUp
		// Đảm bảo có đủ ít nhất 4 phần tử
		if len(row) > 3 && row[3] != nil && row[3].(string) != "" {
			var orgProfile string
			if len(row) > 5 && row[5] != nil && row[5].(string) != "" {
				orgProfile = row[5].(string)
			}
			err = receiver.RoleOrgSignUpRepo.UpdateOrCreate(&entity.SRoleOrgSignUp{
				RoleName:   row[3].(string),
				OrgCode:    code,
				OrgProfile: orgProfile,
			})
			if err != nil {
				return fmt.Errorf("error importing sign up forms step save into RoleOrgSignUp: %s", err.Error())
			}
		}

		url := row[1].(string)
		//validate spreadsheet url
		if !strings.Contains(url, "https://docs.google.com/spreadsheets/d/") {
			continue
		}

		_, err = receiver.importSignUpForm(url, code, sheetName)
		if err != nil {
			return fmt.Errorf("error importing sign up forms: %s", err.Error())
		}
	}

	return nil
}

func (receiver *ImportFormsUseCase) importSignUpForm(spreadsheetUrl, note, sheetNameToRead string) (entity.SForm, error) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(spreadsheetUrl)

	if len(match) < 2 {
		log.Error("Import Sign Up Form Invalid spreadsheet url: ", spreadsheetUrl)
		return entity.SForm{}, fmt.Errorf("import sign up form invalid spreadsheet url: %s", spreadsheetUrl)
	}

	spreadsheetID := match[1]
	monitor.LogGoogleAPIRequestImportForm()
	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetID: spreadsheetID,
		ReadRange:     sheetNameToRead + `!J11` + `:Q`,
	})
	if err != nil {
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
			remember := "false"
			if len(row) > 8 {
				remember = row[8].(string)
			}
			var uniqueID *string = nil
			if row[0].(string) != "" {
				id := row[0].(string)
				uniqueID = &id
			}
			item := parameters.RawQuestion{
				// ID:        strings.ToUpper(note) + "_" + spreadsheetID + "_" + strconv.Itoa(index-1),
				QuestionID:        uuid.NewString(),
				Question:          row[3].(string),
				Type:              row[2].(string),
				Attributes:        strings.ReplaceAll(row[4].(string), "\n", ""),
				AnswerRequired:    required,
				AnswerRemember:    remember,
				AdditionalOptions: additionalInfo,
				Status:            "1",
				RowNumber:         index + 1,
				QuestionUniqueID:  uniqueID,
				Key:               row[0].(string),
				DB:                row[1].(string),
			}
			rawQuestions = append(rawQuestions, item)
		}
	}

	f, msg, err := receiver.CreateSignUpForm(parameters.SaveFormParams{
		Note:           note,
		Name:           formName,
		SpreadsheetUrl: spreadsheetUrl,
		SpreadsheetID:  spreadsheetID,
		Password:       "",
		RawQuestions:   rawQuestions,
		SheetName:      sheetNameToRead,
	})

	if err != nil {
		log.Error(fmt.Sprintf("Error creating form: %s - note : %s", err.Error(), note))
		return entity.SForm{}, err
	}

	log.Warning(msg)

	return *f, nil
}

func (receiver *ImportFormsUseCase) importForm(code string, spreadsheetUrl string, password string, sheetName string) (string, error) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(spreadsheetUrl)

	if len(match) < 2 {
		return "Invalid spreadsheet url", fmt.Errorf("invalid spreadsheet url")
	}

	spreadsheetID := match[1]
	sheetNameToRead := "Questions"
	if sheetName != "" {
		sheetNameToRead = sheetName
	}
	monitor.LogGoogleAPIRequestImportForm()
	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetID: spreadsheetID,
		ReadRange:     sheetNameToRead + `!I11:Q`,
	})
	if err != nil || values == nil {
		log.Error(err)
		return fmt.Sprintf("Cannot read tab %s from %s ERROR %s", sheetNameToRead, spreadsheetUrl, err.Error()), err
	}

	var rawQuestions = make([]parameters.RawQuestion, 0)
	var formName = ""
	for index, row := range values {
		if index == 0 && len(row) > 2 && cap(row) > 2 {
			formName = row[2].(string)
			continue
		} else if len(row) > 5 && cap(row) > 5 && index > 1 && row[3].(string) != "" {
			question := ""
			if len(row) > 4 {
				question = row[4].(string)
			}

			qType := row[3].(string)
			// Rule check: có Type nhưng không có Question
			if len(qType) > 0 && len(strings.TrimSpace(question)) == 0 {
				return fmt.Sprintf("Invalid data at row %d: Question is required when Type is provided", index+1),
					fmt.Errorf("row %d: missing Question while Type = %s", index+1, qType)
			}
			additionalInfo := ""
			if len(row) > 7 {
				additionalInfo = row[7].(string)
			}
			required := "false"
			if len(row) > 6 {
				required = row[6].(string)
			}
			remember := "false"
			if len(row) > 8 {
				remember = row[8].(string)
			}
			var enabled value.QuestionForMobile = value.QuestionForMobile_Enabled
			if strings.ToUpper(row[0].(string)) == "LOCK" {
				enabled = value.QuestionForMobile_Disabled
			}
			item := parameters.RawQuestion{
				// ID:        strings.ToUpper(code) + "_" + spreadsheetID + "_" + row[2].(string),
				QuestionID:        uuid.NewString(),
				Question:          question,
				Type:              qType,
				Attributes:        strings.ReplaceAll(row[5].(string), "\n", ""),
				AnswerRequired:    required,
				AnswerRemember:    remember,
				AdditionalOptions: additionalInfo,
				Status:            "1",
				RowNumber:         index + 1,
				EnableOnMobile:    enabled,
				Key:               row[1].(string),
				DB:                row[2].(string),
			}
			rawQuestions = append(rawQuestions, item)
		}
	}

	_, reason, err := receiver.saveForm(parameters.SaveFormParams{
		Note:           code,
		Name:           formName,
		SpreadsheetUrl: spreadsheetUrl,
		SpreadsheetID:  spreadsheetID,
		Password:       password,
		RawQuestions:   rawQuestions,
		SheetName:      sheetNameToRead,
	})

	return reason, err
}

func (receiver *ImportFormsUseCase) CreateSignUpForm(params parameters.SaveFormParams) (*entity.SForm, string, error) {
	return receiver.saveForm(params)
}

func (receiver *ImportFormsUseCase) saveForm(params parameters.SaveFormParams) (*entity.SForm, string, error) {
	err := receiver.QuestionRepository.DeleteQuestionsFormNote(params.Note)
	if err != nil {
		return nil, "System Error: cannot delete questions belong to this form: " + params.Note, err
	}

	questions, invalidQuestions, err := receiver.saveQuestions(params.RawQuestions)
	var reason string
	if len(invalidQuestions) > 0 {
		reason = "Invalid questions: "
		for _, question := range invalidQuestions {
			reason += "Row No:" + strconv.Itoa(question.RowNumber) + ": " + question.Reason + ", "
		}
	}
	if err != nil {
		return nil, "" + err.Error() + " " + reason, err
	}

	form, err := receiver.createForm(questions, params)

	return form, reason, err
}

type InvalidQuestionRow struct {
	RowNumber int
	Reason    string
}

func (receiver *ImportFormsUseCase) saveQuestions(rawQuestions []parameters.RawQuestion) ([]entity.SQuestion, []InvalidQuestionRow, error) {
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
			ID:               rawQuestion.QuestionID,
			Question:         rawQuestion.Question,
			QuestionType:     strings.ToLower(rawQuestion.Type),
			Attributes:       attString,
			Status:           value.GetRawStatusValue(status),
			Set:              rawQuestion.Attributes,
			EnableOnMobile:   rawQuestion.EnableOnMobile,
			QuestionUniqueID: rawQuestion.QuestionUniqueID,
			Key:              rawQuestion.Key,
			DB:               rawQuestion.DB,
		}
		params = append(params, param)
	}

	if len(params) == 0 {
		//return nil, make([]InvalidQuestionRow, 0), errors.New("this form does not have any valid format questions")
		return make([]entity.SQuestion, 0), make([]InvalidQuestionRow, 0), nil
	}

	questions, err := receiver.QuestionRepository.Create(params)

	return questions, invalidQuestions, err
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
	case value.QuestionTime,
		value.QuestionDate,
		value.QuestionDateTime,
		value.QuestionDurationForward,
		value.QuestionQRCode,
		value.QuestionInText,
		value.QuestionText,
		value.QuestionInCount,
		value.QuestionCount,
		value.QuestionNumber,
		value.QuestionQRCodeFront,
		value.UserInformationValue1,
		value.UserInformationValue2,
		value.UserInformationValue3,
		value.UserInformationValue4,
		value.UserInformationValue5,
		value.UserInformationValue6,
		value.UserInformationValue7,

		value.MessageText1,
		value.MessageText2,
		value.ResponseText1,
		value.ResponseText2,

		value.CameraSquareLens,
		value.OrganizationName,
		value.ApplicationContent,
		value.WaterCup:
		return "{}", nil
	case value.QuestionDurationBackward,
		value.QuestionShowPic,
		value.QuestionButton,
		value.QuestionPlayVideo,
		value.QuestionRandomizer,
		value.QuestionDocument,
		value.QuestionQRCodeGenerator,
		value.QuestionCodeCounting,
		value.QuestionPhoto,
		value.QuestionButtonCount,
		value.QuestionSection,
		value.QuestionFormSection,
		value.QuestionFormSendImmediately,
		value.QuestionSignature,
		value.QuestionSendNotification,

		value.QuestionSignUpPreSetValue1,
		value.QuestionSignUpPreSetValue2,
		value.QuestionSignUpPreSetValue3,

		value.QuestionPresetNickname,
		value.QuestionPresetEmail,
		value.QuestionPresetPassword,
		value.QuestionPresetDob,
		value.QuestionPresetRoleSelectWorkingAddress:
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
	case value.QuestionSelection,
		value.QuestionMultipleChoice,
		value.QuestionSingleChoice,
		value.QuestionChoiceToggle,
		value.QuestionDraggableList,
		value.QuestionPresetConditionAccept,
		value.QuestionPresetRole:
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
	case value.QuestionButtonList:
		re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
		match := re.FindStringSubmatch(rawQuestion.Attributes)

		if len(match) < 2 {
			return "", errors.New("invalid google sheet url")
		}

		return `{"spreadsheet_id" : "` + match[1] + `"}`, nil
	case value.QuestionMessageBox,
		value.PdfViewer,
		value.PdfPicker,
		value.SubmitText,
		value.OutNrTotal,
		value.MemoryText,
		value.OutListEntryHistory,
		value.OutListResponse,
		value.OutNrAverageAll,
		value.OutNrLineGraph,
		value.HiddenMessageText,
		value.Timer,
		value.Chart,
		value.Body:
		message := strings.ReplaceAll(rawQuestion.Attributes, "\n", "\\n")
		jsonMsg := `{"value": "` + message + `"}`
		return jsonMsg, nil
	case value.QuestionWeb:
		attInbase64 := base64.StdEncoding.EncodeToString([]byte(rawQuestion.Attributes))
		return `{"value": "` + attInbase64 + `"}`, nil
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
	case value.SignUpButtonConfiguration1,
		value.SignUpButtonConfiguration2,
		value.SignUpButtonConfiguration3,
		value.SignUpButtonConfiguration4,
		value.SignUpButtonConfiguration5,
		value.SignUpButtonConfiguration6,
		value.SignUpButtonConfiguration7,
		value.SignUpButtonConfiguration8,
		value.SignUpButtonConfiguration9,
		value.SignUpButtonConfiguration10:
		rawButtonConfigurations := strings.Split(rawQuestion.Attributes, ";")
		type ButtonConfiguration struct {
			Color  string `json:"color"`
			QrCode string `json:"qr_code"`
		}
		type ButtonConfigurations struct {
			ButtonConfigurations []ButtonConfiguration `json:"button_configurations"`
		}
		var configurationList = make([]ButtonConfiguration, 0)
		if rawButtonConfigurations[0] == "" || rawButtonConfigurations[1] == "" {
			return "", errors.New("invalid configuration")
		}
		configurationList = append(configurationList, ButtonConfiguration{
			Color:  rawButtonConfigurations[0],
			QrCode: rawButtonConfigurations[1],
		})
		configurations := ButtonConfigurations{ButtonConfigurations: configurationList}
		result, err := json.Marshal(configurations)
		if err != nil {
			return "", err
		}
		return string(result), nil

	default:
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	}
}

func (receiver *ImportFormsUseCase) createForm(questions []entity.SQuestion, params parameters.SaveFormParams) (*entity.SForm, error) {
	form, err := receiver.FormRepository.SaveForm(params)
	if err != nil {
		return nil, err
	}
	formQuestions := make([]request.CreateFormQuestionItem, 0)
	memoryValues := make([]entity.MemoryComponentValue, 0)
	for _, question := range questions {
		var order = 0
		var answerRequired = false
		var answerRemember = false
		for _, rq := range params.RawQuestions {
			if rq.QuestionID == question.ID.String() {
				order = rq.RowNumber
				answerRequired = strings.ToLower(rq.AnswerRequired) == "true"

				if strings.ToLower(rq.AnswerRemember) == "true" {
					answerRemember = strings.ToLower(rq.AnswerRemember) == "true"
					memoryValues = append(memoryValues, entity.MemoryComponentValue{
						ComponentName: question.QuestionType,
					})
				}
			}
		}

		formQuestions = append(formQuestions, request.CreateFormQuestionItem{
			QuestionID:     question.ID.String(),
			Order:          order,
			AnswerRequired: answerRequired,
			AnswerRemember: answerRemember,
		})
	}

	if len(memoryValues) > 0 {
		err = receiver.QuestionRepository.CreateMemoryComponentValues(memoryValues)
	}

	if len(formQuestions) > 0 {
		_, err = receiver.FormQuestionRepository.CreateFormQuestions(form.ID, formQuestions)
		if err != nil {
			return nil, err
		}
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
	spreadsheetID := match[1]
	monitor.LogGoogleAPIRequestImportForm()

	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetID: spreadsheetID,
		ReadRange:     sheetName + `!` + receiver.Google.FirstColumn + strconv.Itoa(receiver.Google.FirstRow+2) + `:AC`,
	})
	if err != nil {
		log.Error(err)
		return fmt.Errorf("error reading spreadsheet: %w", err)
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
						Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.Google.FirstRow+2) + ":Q",
						Dimension: "ROWS",
						Rows:      [][]interface{}{{"DELETED", time.Now().Format("2006-01-02 15:04:05"), ""}},
					}, spreadsheetID)
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
						Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.Google.FirstRow+2) + ":Q",
						Dimension: "ROWS",
						Rows:      [][]interface{}{{"DEACTIVATED", time.Now().Format("2006-01-02 15:04:05"), ""}},
					}, spreadsheetID)
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
				// var submissionType value.SubmissionType = value.SubmissionTypeValues
				// var submissionSheetID string = ""
				// if len(row) >= 15 {
				// 	submissionType = value.GetSubmissionTypeFromString(row[14].(string))
				// 	if len(row) >= 16 {
				// 		//Just required column Z (#15) in case of submission type is in (2,3,5)
				// 		switch submissionType {
				// 		case value.SubmissionTypeBoth, value.SubmissionTypeQrCode, value.SubmissionTypeTeacherAndQRCode: //5
				// 			re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
				// 			match := re.FindStringSubmatch(row[15].(string))

				// 			if len(match) < 2 {
				// 				continue
				// 			} else {
				// 				submissionSheetID = match[1]
				// 			}
				// 		default:
				// 			break
				// 		}
				// 	}
				// }
				// var tabName string = ""
				// if len(row) >= 17 {
				// 	tabName = row[16].(string)
				// }
				// outputSheetName := "Answers"
				// if len(row) >= 18 {
				// 	outputSheetName = row[17].(string)
				// }
				// importErr, reason := receiver.importForm(row[0].(string), row[1].(string), row[2].(string), row[3].(string), submissionType, submissionSheetID, tabName, outputSheetName, syncStrategy)
				// if importErr != nil {
				// 	log.Error(importErr)
				// 	monitor.LogGoogleAPIRequestImportForm()
				// 	_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
				// 		Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
				// 		Dimension: "ROWS",
				// 		Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
				// 	}, spreadsheetID)
				// 	if err != nil {
				// 		log.Debug("Row No: ", rowNo)
				// 		log.Error(err)
				// 	}
				// } else {
				// 	monitor.LogGoogleAPIRequestImportForm()
				// 	_, err = receiver.SpreadsheetWriter.UpdateRange(sheet.WriteRangeParams{
				// 		Range:     sheetName + "!O" + strconv.Itoa(rowNo+receiver.AppConfig.Google.FirstRow+2) + ":Q",
				// 		Dimension: "ROWS",
				// 		Rows:      [][]interface{}{{"UPLOADED", time.Now().Format("2006-01-02 15:04:05"), reason}},
				// 	}, spreadsheetID)
				// 	if err != nil {
				// 		log.Debug("Row No: ", rowNo)
				// 		log.Error(err)
				// 	}
				// }
			}
		}
	}

	return nil
}
