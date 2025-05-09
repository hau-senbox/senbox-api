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
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/sheet"
	"strconv"
	"strings"

	log "github.com/sirupsen/logrus"
)

type SaveFormUseCase struct {
	*repository.FormRepository
	*repository.QuestionRepository
	*repository.FormQuestionRepository
	SpreadsheetReader *sheet.Reader
}

func (receiver *SaveFormUseCase) SaveForm(req request.SaveFormRequest) (*entity.SForm, error) {
	re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
	match := re.FindStringSubmatch(req.SpreadsheetUrl)

	if len(match) < 2 {
		return nil, fmt.Errorf("invalid spreadsheet url")
	}
	spreadsheetId := match[1]
	values, err := receiver.SpreadsheetReader.Get(sheet.ReadSpecificRangeParams{
		SpreadsheetId: spreadsheetId,
		ReadRange:     "Questions!A2:G",
	})
	if err != nil {
		log.Error(err)
		return nil, err
	}

	var rawQuestions = make([]parameters.RawQuestion, 0)
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
			rawQuestions = append(rawQuestions, item)
		}
	}

	form, err := receiver.saveForm(parameters.SaveFormParams{
		Note:           req.Note,
		SpreadsheetUrl: req.SpreadsheetUrl,
		SpreadsheetId:  spreadsheetId,
		Password:       req.Password,
		RawQuestions:   rawQuestions,
	})
	if err != nil {
		return nil, err
	}

	return form, nil
}

func (receiver *SaveFormUseCase) saveForm(params parameters.SaveFormParams) (*entity.SForm, error) {
	questions, err := receiver.saveQuestions(params.RawQuestions)
	if err != nil {
		return nil, err
	}
	log.Debug(questions)

	return receiver.createForm(questions, params)
}

func (receiver *SaveFormUseCase) saveQuestions(rawQuestions []parameters.RawQuestion) ([]entity.SQuestion, error) {
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
	q, err := receiver.QuestionRepository.Create(params)
	return q, err
}

func (receiver *SaveFormUseCase) getStatusFromString(status string) (value.Status, error) {
	switch strings.ToLower(status) {
	case "true":
		return value.Active, nil
	case "false":
		return value.Inactive, nil
	default:
		return value.Active, nil
	}
}

func (receiver *SaveFormUseCase) unmarshalAttributes(rawQuestion parameters.RawQuestion, questionType value.QuestionType) (string, error) {

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
	}

	return "", errors.New("invalid question type")
}

func (receiver *SaveFormUseCase) createForm(questions []entity.SQuestion, params parameters.SaveFormParams) (*entity.SForm, error) {
	form, err := receiver.FormRepository.SaveForm(params)
	if err != nil {
		return nil, err
	}
	var formQuestions = make([]request.CreateFormQuestionItem, 0)
	for _, question := range questions {
		var order = 0
		var answerRequired = false
		for _, rq := range params.RawQuestions {
			if rq.QuestionId == question.QuestionId.String() {
				order = rq.RowNumber
				answerRequired = strings.ToLower(rq.AnswerRequired) == "true"
			}
		}

		formQuestions = append(formQuestions, request.CreateFormQuestionItem{
			QuestionId:     question.QuestionId.String(),
			Order:          order,
			AnswerRequired: answerRequired,
		})
	}
	_, err = receiver.CreateFormQuestions(form.ID, formQuestions)
	if err != nil {
		return nil, err
	}
	return form, nil
}
