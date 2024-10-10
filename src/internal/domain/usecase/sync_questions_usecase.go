package usecase

import (
	"encoding/json"
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/parameters"
	"sen-global-api/internal/domain/value"
	"strings"
)

type SyncQuestionsUseCase struct {
	*repository.QuestionRepository
}

func (receiver *SyncQuestionsUseCase) GetStatusFromString(status string) (value.Status, error) {
	switch strings.ToLower(status) {
	case "true":
		return value.Active, nil
	case "false":
		return value.Inactive, nil
	default:
		return value.Active, nil
	}
}

func (receiver *SyncQuestionsUseCase) unmarshalAttributes(rawQuestion parameters.RawQuestion, questionType value.QuestionType) (string, error) {

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
		return "{}", nil
	case value.QuestionQRCode:
		return "{}", nil
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
	default:
		return "", errors.New("invalid question type")
	}

	return "", errors.New("invalid question type")
}
