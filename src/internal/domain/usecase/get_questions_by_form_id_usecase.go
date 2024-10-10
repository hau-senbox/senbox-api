package usecase

import (
	"encoding/json"
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"net/http"
	"regexp"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"strconv"
	"strings"
)

type GetQuestionsByFormUseCase struct {
	*repository.QuestionRepository
	*repository.DeviceFormDatasetRepository
	*repository.CodeCountingRepository
	*gorm.DB
}

//response.QuestionListResponse
//response.FailedResponse

func (receiver *GetQuestionsByFormUseCase) GetQuestionByForm(form entity.SForm, device entity.SDevice) (*response.QuestionListResponse, *response.FailedResponse) {
	questions, err := receiver.QuestionRepository.GetQuestionsByFormId(form.FormId)
	if err != nil {
		return nil, &response.FailedResponse{
			Error: response.Cause{
				Code:    http.StatusBadRequest,
				Message: "Failed to get questions",
			},
		}
	}
	var result = make([]response.QuestionListData, 0)
	var rawQuestions = make([]response.QuestionListData, 0)
	for _, question := range questions {
		var att response.QuestionAttributes
		err = json.Unmarshal([]byte(question.Attributes), &att)
		if err != nil {
			continue
		}

		// Skip Send Notification
		// because they are not for mobile
		if question.QuestionType == value.GetStringValue(value.QuestionSendNotification) {
			continue
		}

		q := response.QuestionListData{
			QuestionId:     question.QuestionId,
			QuestionType:   strings.ToUpper(question.QuestionType),
			Question:       question.Question,
			Attributes:     att,
			Order:          question.Order,
			AnswerRequired: question.AnswerRequired,
			Enabled:        question.EnableOnMobile == value.QuestionForMobile_Enabled,
		}

		rawQuestions = append(rawQuestions, q)
	}

	for _, rawQuestion := range rawQuestions {
		qType, err := value.GetQuestionType(rawQuestion.QuestionType)
		if err != nil {
			continue
		}
		if value.IsGeneralQuestionType(qType) == true {
			// Check code counting & code generation
			if qType == value.QuestionCodeCounting {
				q, err := receiver.BuildCodeCountingQuestion(rawQuestion)
				if err != nil {
					return nil, &response.FailedResponse{
						Error: response.Cause{
							Code:    555,
							Message: "Could not parsed user form data fo question: " + rawQuestion.QuestionId + " err: " + err.Error(),
						},
					}
				}
				result = append(result, q)
			} else if qType == value.QuestionRandomizer {
				q, err := receiver.BuildRandomizerQuestion(rawQuestion)
				if err != nil {
					return nil, &response.FailedResponse{
						Error: response.Cause{
							Code:    555,
							Message: "Could not parsed user form data fo question: " + rawQuestion.QuestionId + " err: " + err.Error(),
						},
					}
				}
				result = append(result, q)
			} else {
				result = append(result, rawQuestion)
			}
		} else {
			q, err := receiver.BuildQuestion(rawQuestion, qType, device.DeviceId)
			if err != nil {
				return nil, &response.FailedResponse{
					Error: response.Cause{
						Code:    555,
						Message: "Could not parsed user form data fo question: " + rawQuestion.QuestionId + " err: " + err.Error(),
					},
				}
			}
			if q == nil {
				return nil, &response.FailedResponse{
					Error: response.Cause{
						Code:    555,
						Message: "Could not parsed user form data fo question: " + rawQuestion.QuestionId,
					},
				}
			}
			result = append(result, *q)
		}
	}

	return &response.QuestionListResponse{
		Data: response.QuestionListResponseData{
			QuestionListData: result,
			DecryptPassword:  form.Password,
			FormName:         form.Name,
		},
	}, nil
}

func (receiver *GetQuestionsByFormUseCase) BuildQuestion(question response.QuestionListData, qType value.QuestionType, deviceId string) (*response.QuestionListData, error) {
	dataset, err := receiver.DeviceFormDatasetRepository.GetDatasetByDeviceIdAndSet(deviceId, question.Attributes.Value)
	if err != nil {
		return nil, err
	}
	var att response.QuestionAttributes
	attInJSONString := "{}"
	switch qType {
	case value.QuestionDateUser:
		attInJSONString = "{}"
	case value.QuestionTimeUser:
		attInJSONString = "{}"
	case value.QuestionDateTimeUser:
		attInJSONString = "{}"
	case value.QuestionDurationForwardUser:
		attInJSONString = "{}"
	case value.QuestionDurationBackwardUser:
		attInJSONString = `{"value": "` + dataset.QuestionDurationBackward + `"}`
	case value.QuestionScaleUser:
		rawValues := strings.Split(dataset.QuestionScale, ";")
		if len(rawValues) < 2 {
			return nil, errors.New("scale question data is invalid " + dataset.QuestionScale)
		}
		totalValuesInString := strings.Split(rawValues[0], ":")
		if len(totalValuesInString) < 2 {
			return nil, errors.New("scale question data is invalid " + dataset.QuestionScale)
		}
		stepValueInString := strings.Split(rawValues[1], ":")
		if len(stepValueInString) < 2 {
			return nil, errors.New("scale question data is invalid " + dataset.QuestionScale)
		}
		totalValues, err := strconv.Atoi(totalValuesInString[1])
		if err != nil {
			return nil, errors.New("scale question data is invalid " + err.Error())
		}
		stepValue, err := strconv.Atoi(stepValueInString[1])
		if err != nil {
			return nil, errors.New("scale question data is invalid " + err.Error())
		}

		attInJSONString = "{\"number\" : " + strconv.Itoa(totalValues) + ", \"steps\": " + strconv.Itoa(stepValue) + "}"
	case value.QuestionQRCodeUser:
		attInJSONString = "{}"
	case value.QuestionSelectionUser:
		rawOptions := strings.Split(dataset.QuestionSelection, ";")
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
				return nil, errors.New("invalid options")
			}
			optionsList = append(optionsList, Option{Name: op})
		}
		options := Options{Options: optionsList}
		result, err := json.Marshal(options)
		if err != nil {
			return nil, err
		}
		attInJSONString = string(result)
	case value.QuestionTextUser:
		attInJSONString = "{}"
	case value.QuestionCountUser:
		attInJSONString = "{}"
	case value.QuestionNumberUser:
		attInJSONString = "{}"
	case value.QuestionPhotoUser:
		attInJSONString = "{}"
	case value.QuestionButtonCountUser:
		attInJSONString = `{"value": "` + dataset.QuestionButtonCount + `"}`
	case value.QuestionMultipleChoiceUser:
		rawOptions := strings.Split(dataset.QuestionMultipleChoice, ";")
		type Option struct {
			Name string `json:"name"`
		}
		type Options struct {
			Options []Option `json:"options"`
		}
		var optionsList = make([]Option, 0)
		for _, op := range rawOptions {
			if op == "" {
				return nil, errors.New("invalid options")
			}
			optionsList = append(optionsList, Option{Name: op})
		}
		options := Options{Options: optionsList}
		result, err := json.Marshal(options)
		if err != nil {
			return nil, err
		}
		attInJSONString = string(result)
	case value.QuestionSingleChoiceUser:
		rawOptions := strings.Split(dataset.QuestionSingleChoice, ";")
		type Option struct {
			Name string `json:"name"`
		}
		type Options struct {
			Options []Option `json:"options"`
		}
		var optionsList = make([]Option, 0)
		for _, op := range rawOptions {
			if op == "" {
				return nil, errors.New("invalid options")
			}
			optionsList = append(optionsList, Option{Name: op})
		}
		options := Options{Options: optionsList}
		result, err := json.Marshal(options)
		if err != nil {
			return nil, err
		}
		attInJSONString = string(result)
	case value.QuestionButtonListUser:
		re := regexp.MustCompile(`/spreadsheets/d/([a-zA-Z0-9-_]+)`)
		match := re.FindStringSubmatch(dataset.QuestionButtonList)

		if len(match) < 2 {
			return nil, errors.New("invalid google sheet url")
		}

		attInJSONString = `{"spreadsheet_id" : "` + match[1] + `"}`

	case value.QuestionMessageBoxUser:
		message := strings.Replace(dataset.QuestionMessageBox, "\n", "\\n", -1)
		jsonMsg := `{"value": "` + message + `"}`
		attInJSONString = jsonMsg
	case value.QuestionShowPicUser:
		attInJSONString = `{"value": "` + dataset.QuestionShowPic + `"}`
	case value.QuestionButtonUser:
		attInJSONString = `{"value": "` + dataset.QuestionButton + `"}`
	case value.QuestionPlayVideoUser:
		attInJSONString = `{"value": "` + dataset.QuestionPlayVideo + `"}`
	case value.QuestionQRCodeFrontUser:
		attInJSONString = "{}"
	case value.QuestionChoiceToggleUser:
		rawOptions := strings.Split(dataset.QuestionChoiceToggle, ";")
		type Option struct {
			Name string `json:"name"`
		}
		type Options struct {
			Options []Option `json:"options"`
		}
		var optionsList = make([]Option, 0)
		for _, op := range rawOptions {
			if op == "" {
				return nil, errors.New("invalid options")
			}
			optionsList = append(optionsList, Option{Name: op})
		}
		options := Options{Options: optionsList}
		result, err := json.Marshal(options)
		if err != nil {
			return nil, err
		}
		attInJSONString = string(result)
	case value.QuestionSignature:
		attInJSONString = `{"value": "` + dataset.QuestionSignature + `"}`
	case value.QuestionWeb:
		attInJSONString = `{"value": "` + dataset.QuestionWeb + `"}`
	case value.QuestionWebUser:
		attInJSONString = `{"value": "` + dataset.QuestionWeb + `"}`

	default:
		return nil, errors.New("invalid question type for user form")
	}

	err = json.Unmarshal([]byte(attInJSONString), &att)
	if err != nil {
		return nil, err
	}

	q := response.QuestionListData{
		QuestionId:     question.QuestionId,
		QuestionType:   strings.ToUpper(question.QuestionType),
		Question:       question.Question,
		Attributes:     att,
		Order:          question.Order,
		AnswerRequired: question.AnswerRequired,
		Enabled:        question.Enabled,
	}

	return &q, nil
}

func (receiver *GetQuestionsByFormUseCase) GetQuestionsBySignUpForm(form entity.SForm) *response.QuestionListResponse {
	questions, err := receiver.QuestionRepository.GetQuestionsByFormId(form.FormId)
	if err != nil {
		return nil
	}
	var result = make([]response.QuestionListData, 0)
	var rawQuestions = make([]response.QuestionListData, 0)
	for _, question := range questions {
		var att response.QuestionAttributes
		err = json.Unmarshal([]byte(question.Attributes), &att)
		if err != nil {
			continue
		}
		q := response.QuestionListData{
			QuestionId:     question.QuestionId,
			QuestionType:   strings.ToUpper(question.QuestionType),
			Question:       question.Question,
			Attributes:     att,
			Order:          question.Order,
			AnswerRequired: question.AnswerRequired,
			Enabled:        question.EnableOnMobile == value.QuestionForMobile_Enabled,
		}

		rawQuestions = append(rawQuestions, q)
	}

	for _, rawQuestion := range rawQuestions {
		qType, err := value.GetQuestionType(rawQuestion.QuestionType)
		if err != nil {
			continue
		}
		if value.IsGeneralQuestionType(qType) == true {
			result = append(result, rawQuestion)
		}
	}

	return &response.QuestionListResponse{
		Data: response.QuestionListResponseData{
			QuestionListData: result,
			DecryptPassword:  form.Password,
			FormName:         form.Name,
		},
	}
}

func (receiver *GetQuestionsByFormUseCase) BuildCodeCountingQuestion(question response.QuestionListData) (response.QuestionListData, error) {
	var att response.QuestionAttributes
	attInJSONString := "{}"

	newCodeCountingValue, err := receiver.CodeCountingRepository.CreateForQuestionWithID(question.QuestionId, receiver.DB)
	if err != nil {
		log.Error(err)
		return response.QuestionListData{}, err
	}
	attInJSONString = `{"value": "` + newCodeCountingValue + `"}`

	err = json.Unmarshal([]byte(attInJSONString), &att)
	if err != nil {
		return response.QuestionListData{}, err
	}

	q := response.QuestionListData{
		QuestionId:     question.QuestionId,
		QuestionType:   strings.ToUpper(question.QuestionType),
		Question:       question.Question,
		Attributes:     att,
		Order:          question.Order,
		AnswerRequired: question.AnswerRequired,
		Enabled:        question.Enabled,
	}

	return q, nil
}

func (receiver *GetQuestionsByFormUseCase) BuildRandomizerQuestion(question response.QuestionListData) (response.QuestionListData, error) {
	var att response.QuestionAttributes
	attInJSONString := "{}"

	code := att.Value + value.GetRandomString(8)
	attInJSONString = `{"value": "` + question.Attributes.Value + code + `"}`

	err := json.Unmarshal([]byte(attInJSONString), &att)
	if err != nil {
		return response.QuestionListData{}, err
	}

	q := response.QuestionListData{
		QuestionId:     question.QuestionId,
		QuestionType:   strings.ToUpper(question.QuestionType),
		Question:       question.Question,
		Attributes:     att,
		Order:          question.Order,
		AnswerRequired: question.AnswerRequired,
		Enabled:        question.Enabled,
	}

	return q, nil
}
