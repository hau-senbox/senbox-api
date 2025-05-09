package usecase

import (
	"encoding/json"
	"errors"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/parameters"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type UpdateFormUseCase struct {
	*repository.FormRepository
	*repository.QuestionRepository
	*repository.FormQuestionRepository
	SpreadsheetReader *sheet.Reader
}

func (receiver *UpdateFormUseCase) UpdateForm(formId int, request request.UpdateFormRequest) (*entity.SForm, error) {
	form, err := receiver.GetFormById(uint64(formId))
	if err != nil {
		return nil, err
	}
	if request.Password != nil {
		form.Password = *request.Password
	}

	rawQuestions, err := receiver.getRawQuestions(form.SpreadsheetId)
	if err != nil {
		return nil, err
	}
	questions, err := receiver.syncQuestions(rawQuestions)
	if err != nil {
		return nil, err
	}

	err = receiver.Update(form, questions, rawQuestions)
	if err != nil {
		return nil, err
	}

	log.Debug(questions)

	return receiver.FormRepository.UpdateForm(form)
}

func (receiver *UpdateFormUseCase) getRawQuestions(spreadsheetId string) ([]parameters.RawQuestion, error) {
	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Questions!A2:G",
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var result = make([]parameters.RawQuestion, 0)
	for index, row := range values {
		if len(row) >= 7 && cap(row) >= 7 {
			if row[0].(string) == "" {
				continue
			}
			item := parameters.RawQuestion{
				QuestionId:        spreadsheetId + "_" + row[0].(string),
				Question:          row[2].(string),
				Type:              row[1].(string),
				Attributes:        strings.ReplaceAll(row[3].(string), "\n", ""),
				AnswerRequired:    row[4].(string),
				AdditionalOptions: row[5].(string),
				Status:            row[6].(string),
				RowNumber:         index + 1,
			}
			result = append(result, item)
		}
	}

	return result, err
}

func (receiver *UpdateFormUseCase) syncQuestions(rawQuestions []parameters.RawQuestion) ([]entity.SQuestion, error) {
	var params = make([]repository.CreateQuestionParams, 0)
	for _, rawQuestion := range rawQuestions {
		questionType, err := value.GetQuestionType(rawQuestion.Type)
		if err != nil {
			continue
		}

		status, err := receiver.getStatusFromString(rawQuestion.Status)
		if err != nil {
			continue
		}

		attString, err := receiver.unmarshalAttributes(rawQuestion, questionType)
		if err != nil {
			continue
		}

		param := repository.CreateQuestionParams{
			QuestionId:   rawQuestion.QuestionId,
			QuestionName: rawQuestion.Question,
			QuestionType: strings.ToLower(rawQuestion.Type),
			Question:     rawQuestion.Question,
			Attributes:   attString,
			Status:       value.GetRawStatusValue(status),
		}
		params = append(params, param)
	}
	return receiver.QuestionRepository.Create(params)
}

func (receiver *UpdateFormUseCase) getStatusFromString(status string) (value.Status, error) {
	switch strings.ToLower(status) {
	case "true":
		return value.Active, nil
	case "false":
		return value.Inactive, nil
	default:
		return value.Active, nil
	}
}

func (receiver *UpdateFormUseCase) unmarshalAttributes(rawQuestion parameters.RawQuestion, questionType value.QuestionType) (string, error) {

	switch questionType {
	case value.QuestionTime,
		value.QuestionDate,
		value.QuestionDateTime,
		value.QuestionDurationForward,
		value.QuestionQRCode,
		value.QuestionText,
		value.QuestionCount,
		value.QuestionNumber:
		return "{}", nil
	case value.QuestionDurationBackward,
		value.QuestionPhoto,
		value.QuestionButtonCount,
		value.QuestionMessageBox,
		value.QuestionShowPic,
		value.QuestionButton,
		value.QuestionPlayVideo:
		//TODO: validate attributes
		return `{"value": "` + rawQuestion.Attributes + `"}`, nil
	case value.QuestionScale:
		rawValues := strings.Split(rawQuestion.Attributes, ",")
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
	case value.QuestionSelection:
		rawOptions := strings.Split(rawQuestion.Attributes, ",")
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
	case value.QuestionMultipleChoice:
		multiselect := "single_select"
		rawAdditionalOptions := strings.Split(rawQuestion.AdditionalOptions, ":")
		if len(rawAdditionalOptions) > 1 {
			if strings.ToLower(strings.TrimSpace(rawAdditionalOptions[1])) == "yes" {
				multiselect = "multi_select"
			}
		}
		rawOptions := strings.Split(rawQuestion.Attributes, ",")
		type Option struct {
			Name string `json:"name"`
		}
		type Options struct {
			Options []Option `json:"options"`
			Value   string   `json:"value"`
		}
		var optionsList = make([]Option, 0)
		for _, op := range rawOptions {
			if op == "" {
				return "", errors.New("invalid options")
			}
			optionsList = append(optionsList, Option{Name: op})
		}
		options := Options{Options: optionsList, Value: multiselect}
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
	default:
		return "", errors.New("invalid question type")
	}
}
