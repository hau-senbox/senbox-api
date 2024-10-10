package usecase

import (
	"encoding/json"
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	log "github.com/sirupsen/logrus"
)

type DeviceSignUpTextButtonSetting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type DeviceSignUpFormSetting struct {
	FormId uint64 `json:"form_id"`
}

type DeviceSignUpSetting struct {
	Button1 DeviceSignUpTextButtonSetting
	Button2 DeviceSignUpTextButtonSetting
	Button3 DeviceSignUpTextButtonSetting
	Button4 DeviceSignUpTextButtonSetting
}

type DeviceSignUpUseCases struct {
	SettingRepository         *repository.SettingRepository
	FormRepository            *repository.FormRepository
	GetQuestionsByFormUseCase *GetQuestionsByFormUseCase
}

func (c *DeviceSignUpUseCases) GetSignUpSetting() (DeviceSignUpSetting, error) {
	btt1, err := c.SettingRepository.GetSignUpButton1Setting()
	if err != nil {
		return DeviceSignUpSetting{}, err
	}

	btt2, err := c.SettingRepository.GetSignUpButton2Setting()
	if err != nil {
		return DeviceSignUpSetting{}, err
	}

	btt3, err := c.SettingRepository.GetSignUpButton3Setting()
	if err != nil {
		return DeviceSignUpSetting{}, err
	}

	btt4, err := c.SettingRepository.GetSignUpButton4Setting()
	if err != nil {
		return DeviceSignUpSetting{}, err
	}

	formUrl, err := c.SettingRepository.GetRegistrationFormSetting()
	if err != nil {
		return DeviceSignUpSetting{}, err
	}

	if btt1 == nil || btt2 == nil || btt3 == nil || btt4 == nil || formUrl == nil {
		return DeviceSignUpSetting{}, errors.New("Some SignUp settings are not set")
	}

	var btt1Setting DeviceSignUpTextButtonSetting
	if err := json.Unmarshal(btt1.Settings, &btt1Setting); err != nil {
		return DeviceSignUpSetting{}, err
	}

	var btt2Setting DeviceSignUpTextButtonSetting
	if err := json.Unmarshal(btt2.Settings, &btt2Setting); err != nil {
		return DeviceSignUpSetting{}, err
	}

	var btt3Setting DeviceSignUpTextButtonSetting
	if err := json.Unmarshal(btt3.Settings, &btt3Setting); err != nil {
		return DeviceSignUpSetting{}, err
	}

	var btt4Setting DeviceSignUpTextButtonSetting
	if err := json.Unmarshal(btt4.Settings, &btt4Setting); err != nil {
		return DeviceSignUpSetting{}, err
	}

	var formUrlSetting DeviceSignUpFormSetting
	if err := json.Unmarshal(formUrl.Settings, &formUrlSetting); err != nil {
		return DeviceSignUpSetting{}, err
	}

	return DeviceSignUpSetting{
		Button1: btt1Setting,
		Button2: btt2Setting,
		Button3: btt3Setting,
		Button4: btt4Setting,
	}, nil
}

func (c *DeviceSignUpUseCases) GetSignUpFormQuestions() *response.QuestionListResponse {
	f, err := c.findSignUpForm()
	if err != nil {
		return nil
	}

	return c.GetQuestionsByFormUseCase.GetQuestionsBySignUpForm(f)
}

func (c *DeviceSignUpUseCases) GetSigGnUpPresetSetting() *string {
	s, err := c.SettingRepository.GetRegistrationPresetSetting()

	if err != nil || s == nil {
		log.Error(err)
	}

	type SummarySetting struct {
		SpreadsheetId string `json:"spreadsheet_id"`
	}

	var registrationPresetSettings *SummarySetting = nil

	err = json.Unmarshal([]byte(s.Settings), &registrationPresetSettings)
	if err != nil {
		log.Info(err.Error())

		return nil
	}

	return &registrationPresetSettings.SpreadsheetId
}

func (c *DeviceSignUpUseCases) findSignUpForm() (entity.SForm, error) {
	s, err := c.SettingRepository.GetRegistrationFormSetting()
	if err != nil {
		return entity.SForm{}, err
	}

	if s == nil {
		return entity.SForm{}, errors.New("Sign Up Form not found")
	}

	//Unmarshal
	var formSetting DeviceSignUpFormSetting
	if err := json.Unmarshal(s.Settings, &formSetting); err != nil {
		return entity.SForm{}, err
	}

	f, err := c.FormRepository.GetFormById(formSetting.FormId)
	if err != nil {
		return entity.SForm{}, err
	}

	if f == nil {
		return entity.SForm{}, errors.New("Form not found")
	}

	return *f, nil
}

func (c *DeviceSignUpUseCases) FindSignUpForm() (entity.SForm, error) {
	return c.findSignUpForm()
}

func (c *DeviceSignUpUseCases) GetSignUpFormQuestionsByDevice(deviceID string) *response.QuestionListResponse {
	f, err := c.FormRepository.FindSignUpFormByDeviceID(deviceID)

	if err != nil {
		return nil
	}

	return c.GetQuestionsByFormUseCase.GetQuestionsBySignUpForm(f)
}

func (c *DeviceSignUpUseCases) GetSignUpFormQuestionsByFormNote(code string) *response.QuestionListResponse {
	f, err := c.findSignUpFormByNote(code)

	if err != nil {
		return nil
	}

	return c.GetQuestionsByFormUseCase.GetQuestionsBySignUpForm(f)
}

func (c *DeviceSignUpUseCases) findSignUpFormByNote(code string) (entity.SForm, error) {
	f, err := c.FormRepository.FindByCode(code)

	if err != nil {
		return entity.SForm{}, err
	}

	if f.SubmissionType != value.SubmissionTypeSignUpRegistration {
		return entity.SForm{}, errors.New("Form not found")
	}

	return f, nil
}
