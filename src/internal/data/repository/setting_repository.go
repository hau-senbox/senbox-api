package repository

import (
	"encoding/json"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SettingRepository struct {
	DBConn *gorm.DB
}

func NewSettingRepository(dbConn *gorm.DB) *SettingRepository {
	return &SettingRepository{DBConn: dbConn}
}

func (receiver *SettingRepository) GetFormSettings() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeImportForms)
}

func (receiver *SettingRepository) GetFormSettings2() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeImportForms2)
}

func (receiver *SettingRepository) GetFormSettings3() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeImportForms3)
}

func (receiver *SettingRepository) GetFormSettings4() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeImportForms4)
}

func (receiver *SettingRepository) GetUrlSettings() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeImportUrls)
}

func (receiver *SettingRepository) GetOutputSettings() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSubmission)
}

func (receiver *SettingRepository) GetSummarySettings() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSummary)
}

func (receiver *SettingRepository) getSettingsByType(settingType value.SettingType) (*entity.SSetting, error) {
	var setting entity.SSetting
	err := receiver.DBConn.Where("type = ?", settingType).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

type ImportSetting struct {
	SpreadSheetUrl string `json:"spreadsheet_url"`
	AutoImport     bool   `json:"auto"`
	Interval       uint64 `json:"interval"`
}

type OutputSetting struct {
	FolderId string `json:"folder_id"`
}

type SummarySetting struct {
	SpreadSheetId string `json:"spreadsheet_id"`
}

type SignUpTextButtonSetting struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type SignUpFormSetting struct {
	FormId        uint64 `json:"form_id"`
	SpreadSheetId string `json:"spreadsheet_id"`
}

func (receiver *SettingRepository) UpdateFormSetting(req request.ImportFormRequest) error {
	return receiver.updateFormSettingByType(req, value.SettingTypeImportForms)
}

func (receiver *SettingRepository) UpdateFormSetting2(req request.ImportFormRequest) error {
	return receiver.updateFormSettingByType(req, value.SettingTypeImportForms2)
}

func (receiver *SettingRepository) UpdateFormSetting3(req request.ImportFormRequest) error {
	return receiver.updateFormSettingByType(req, value.SettingTypeImportForms3)
}

func (receiver *SettingRepository) UpdateFormSetting4(req request.ImportFormRequest) error {
	return receiver.updateFormSettingByType(req, value.SettingTypeImportForms4)
}

func (receiver *SettingRepository) updateFormSettingByType(req request.ImportFormRequest, formUploaderType value.SettingType) error {
	setting, err := receiver.getSettingsByType(formUploaderType)
	if setting == nil {
		importSetting := ImportSetting{
			SpreadSheetUrl: req.SpreadsheetUrl,
			AutoImport:     req.AutoImport,
			Interval:       req.Interval,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     formUploaderType,
		})

		return result.Error
	} else {
		importSetting := ImportSetting{
			SpreadSheetUrl: req.SpreadsheetUrl,
			AutoImport:     req.AutoImport,
			Interval:       req.Interval,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return err
		}
	}

	return err
}

func (receiver *SettingRepository) UpdateUrlSetting(req request.ImportRedirectUrlsRequest) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeImportUrls)
	if setting == nil {
		importSetting := ImportSetting{
			SpreadSheetUrl: req.SpreadsheetUrl,
			AutoImport:     req.AutoImport,
			Interval:       req.Interval,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeImportUrls,
		})

		return result.Error
	} else {
		importSetting := ImportSetting{
			SpreadSheetUrl: req.SpreadsheetUrl,
			AutoImport:     req.AutoImport,
			Interval:       req.Interval,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return err
		}
	}

	return err
}

func (receiver *SettingRepository) UpdateSubmissionSetting(id string, name string) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeSubmission)
	if setting == nil {
		importSetting := OutputSetting{
			FolderId: id,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSubmission,
		})

		return result.Error
	} else {
		importSetting := OutputSetting{
			FolderId: id,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return err
		}
	}

	return err
}

func (receiver *SettingRepository) GetSubmissionSetting() (OutputSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSubmission)
	if err != nil {
		return OutputSetting{}, err
	}
	var outputSetting OutputSetting
	err = json.Unmarshal([]byte(setting.Settings), &outputSetting)
	if err != nil {
		return OutputSetting{}, err
	}
	return outputSetting, nil
}

func (receiver *SettingRepository) UpdateOutputSummarySetting(spreadsheetId string) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeSummary)
	if setting == nil {
		importSetting := SummarySetting{
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSummary,
		})

		return result.Error
	} else {
		importSetting := SummarySetting{
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return err
		}
	}

	return err
}

func (receiver *SettingRepository) GetSyncDevicesSettings() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSyncDevices)
}

func (receiver *SettingRepository) GetSyncToDosSettings() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSyncToDos)
}

func (receiver *SettingRepository) UpdateSyncToDoSetting(req request.ImportFormRequest) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeSyncToDos)
	if setting == nil {
		importSetting := ImportSetting{
			SpreadSheetUrl: req.SpreadsheetUrl,
			AutoImport:     req.AutoImport,
			Interval:       req.Interval,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSyncToDos,
		})

		return result.Error
	} else {
		importSetting := ImportSetting{
			SpreadSheetUrl: req.SpreadsheetUrl,
			AutoImport:     req.AutoImport,
			Interval:       req.Interval,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return err
		}
	}

	return err
}

func (receiver *SettingRepository) GetEmailSettings() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeEmailHistory)
}

func (receiver *SettingRepository) UpdateEmaiHistorySetting(spreadsheetId string) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeEmailHistory)
	if setting == nil {
		importSetting := SummarySetting{
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeEmailHistory,
		})

		return result.Error
	} else {
		importSetting := SummarySetting{
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return err
		}
	}

	return err
}

func (receiver *SettingRepository) UpdateOutputTemplateSetting(spreadsheetId string) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeOutputTemplate)
	if setting == nil {
		importSetting := SummarySetting{
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeOutputTemplate,
		})

		return result.Error
	} else {
		importSetting := SummarySetting{
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return err
		}
	}

	return err
}

func (receiver *SettingRepository) GetOutputTemplateSettings() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeOutputTemplate)
}

func (receiver *SettingRepository) GetOutputTemplateSettingsForTeacher() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeOutputTemplateTeacher)
}

func (receiver *SettingRepository) UpdateOutputTemplateSettingForTeacher(spreadsheetId string) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeOutputTemplateTeacher)
	if setting == nil {
		importSetting := SummarySetting{
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeOutputTemplateTeacher,
		})

		return result.Error
	} else {
		importSetting := SummarySetting{
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return err
		}
	}

	return err
}

func (receiver *SettingRepository) GetSignUpButton1Setting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSignUpButton1)
}

func (receiver *SettingRepository) GetSignUpButton2Setting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSignUpButton2)
}

func (receiver *SettingRepository) GetSignUpButton3Setting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSignUpButton3)
}

func (receiver *SettingRepository) GetSignUpButton4Setting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSignUpButton4)
}

func (receiver *SettingRepository) GetSignUpButton5Setting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSignUpButton5)
}

func (receiver *SettingRepository) GetSignUpButtonConfigurationSetting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSignUpButtonConfiguration)
}

func (receiver *SettingRepository) GetRegistrationFormSetting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSignUpForm)
}

func (receiver *SettingRepository) GetRegistrationSubmissionSetting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSignUpOutput)
}

func (receiver *SettingRepository) GetRegistrationPreset2Setting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSignUpPresetValue2)
}

func (receiver *SettingRepository) UpdateSignUpButton1(name string, v string) (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSignUpButton1)
	if setting == nil {
		importSetting := SignUpTextButtonSetting{
			Name:  name,
			Value: v,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSignUpButton1,
		})

		return nil, result.Error
	} else {
		importSetting := SignUpTextButtonSetting{
			Name:  name,
			Value: v,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return nil, err
		}
	}

	return setting, err
}

func (receiver *SettingRepository) UpdateSignUpButton2(name string, v string) (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSignUpButton2)
	if setting == nil {
		importSetting := SignUpTextButtonSetting{
			Name:  name,
			Value: v,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSignUpButton2,
		})

		return nil, result.Error
	} else {
		importSetting := SignUpTextButtonSetting{
			Name:  name,
			Value: v,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return nil, err
		}
	}

	return setting, err
}

func (receiver *SettingRepository) UpdateSignUpButton3(name string, v string) (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSignUpButton3)
	if setting == nil {
		importSetting := SignUpTextButtonSetting{
			Name:  name,
			Value: v,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSignUpButton3,
		})

		return nil, result.Error
	} else {
		importSetting := SignUpTextButtonSetting{
			Name:  name,
			Value: v,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return nil, err
		}
	}

	return setting, err
}

func (receiver *SettingRepository) UpdateSignUpButton4(name string, v string) (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSignUpButton4)
	if setting == nil {
		importSetting := SignUpTextButtonSetting{
			Name:  name,
			Value: v,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSignUpButton4,
		})

		return nil, result.Error
	} else {
		importSetting := SignUpTextButtonSetting{
			Name:  name,
			Value: v,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return nil, err
		}
	}

	return setting, err
}

func (receiver *SettingRepository) UpdateSignUpButton5(name string, v string) (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSignUpButton5)
	if setting == nil {
		importSetting := SignUpTextButtonSetting{
			Name:  name,
			Value: v,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSignUpButton5,
		})

		return nil, result.Error
	} else {
		importSetting := SignUpTextButtonSetting{
			Name:  name,
			Value: v,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return nil, err
		}
	}

	return setting, err
}

func (receiver *SettingRepository) UpdateRegistrationForm(formId uint64, url string) (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSignUpForm)
	if setting == nil {
		importSetting := SignUpFormSetting{
			FormId:        formId,
			SpreadSheetId: url,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSignUpForm,
		})

		return nil, result.Error
	} else {
		importSetting := SignUpFormSetting{
			FormId:        formId,
			SpreadSheetId: url,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return nil, err
		}
	}
	return setting, err
}

func (receiver *SettingRepository) UpdateSignUpButtonConfiguration(formId uint64, url string) (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSignUpButtonConfiguration)
	type SignUpButton struct {
		FormId        uint64 `json:"form_id"`
		SpreadSheetId string `json:"spreadsheet_id"`
	}
	if setting == nil {
		importSetting := SignUpButton{
			FormId:        formId,
			SpreadSheetId: url,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSignUpButtonConfiguration,
		})

		return nil, result.Error
	} else {
		importSetting := SignUpButton{
			FormId:        formId,
			SpreadSheetId: url,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return nil, err
		}
	}
	return setting, err
}

func (receiver *SettingRepository) UpdateRegistrationSubmission(url string) (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSignUpOutput)
	if setting == nil {
		importSetting := SummarySetting{
			SpreadSheetId: url,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSignUpOutput,
		})

		return nil, result.Error
	} else {
		importSetting := SummarySetting{
			SpreadSheetId: url,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return nil, err
		}
	}
	return setting, err
}

func (receiver *SettingRepository) UpdateRegistrationPreset2(url string) (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSignUpPresetValue2)
	if setting == nil {
		importSetting := SignUpFormSetting{
			SpreadSheetId: url,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSignUpPresetValue2,
		})

		return nil, result.Error
	} else {
		importSetting := SignUpFormSetting{
			SpreadSheetId: url,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return nil, err
		}
	}

	return setting, err
}

func (receiver *SettingRepository) SetName(name string, settingType value.SettingType) error {
	return receiver.DBConn.
		Model(&entity.SSetting{}).
		Where("type = ?", settingType).
		Update("setting_name", name).
		Error
}

type APIDistributorSetting struct {
	Url           string `json:"spreadsheet_url"`
	SpreadSheetId string `json:"spreadsheet_id"`
}

func (receiver *SettingRepository) UpdateAPIDistributerSetting(spreadsheetId string, url string) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeAPIDistributer)
	if setting == nil {
		importSetting := APIDistributorSetting{
			Url:           url,
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeAPIDistributer,
		})

		return result.Error
	} else {
		importSetting := APIDistributorSetting{
			Url:           url,
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return err
		}
	}

	return err
}

func (receiver *SettingRepository) GetAPIDistributerSetting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeAPIDistributer)
}

func (receiver *SettingRepository) GetCodeCountingDataSetting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeCodeCountingData)
}

func (receiver *SettingRepository) UpdateCodeCountingDataSetting(spreadsheetId string, url string) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeCodeCountingData)
	if setting == nil {
		importSetting := APIDistributorSetting{
			Url:           url,
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: b,
			Type:     value.SettingTypeCodeCountingData,
		})

		return result.Error
	} else {
		importSetting := APIDistributorSetting{
			Url:           url,
			SpreadSheetId: spreadsheetId,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return err
		}
		setting.Settings = b

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return err
		}
	}

	return err
}

func (receiver *SettingRepository) UpdateLogoRefreshInterval(interval uint64) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeLogoRefreshInterval)
	if err != nil {
		if setting == nil {
			result := receiver.DBConn.Create(&entity.SSetting{
				Type:         value.SettingTypeLogoRefreshInterval,
				IntegerValue: interval,
			})
			return result.Error
		} else {
			result := receiver.DBConn.Model(&entity.SSetting{}).Where("type = ?", value.SettingTypeLogoRefreshInterval).Update("settings", interval)
			return result.Error
		}
	} else {
		setting.IntegerValue = interval
		err = receiver.DBConn.Save(&setting).Error
	}

	return err
}

func (receiver *SettingRepository) UpdateLogoRefreshTitle(title string) error {
	setting, err := receiver.getSettingsByType(value.SettingTypeLogoRefreshInterval)
	if err != nil {
		err = receiver.DBConn.Create(&entity.SSetting{
			Type:         value.SettingTypeLogoRefreshInterval,
			SettingName:  title,
			IntegerValue: 10,
		}).Error
	} else {
		setting.SettingName = title
		err = receiver.DBConn.Save(&setting).Error
	}

	return err
}

func (receiver *SettingRepository) GetLogoRefreshInterval() (entity.SSetting, error) {
	var setting entity.SSetting
	err := receiver.DBConn.Where("type = ?", value.SettingTypeLogoRefreshInterval).First(&setting).Error
	if err != nil {
		return setting, err
	}

	return setting, nil
}

func (receiver *SettingRepository) GetImportSignUpFormsSetting() (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeImportSignUpForms)
	if err != nil {
		return nil, err
	}

	return setting, nil
}

func (receiver *SettingRepository) UpdateSignUpFormSetting(req request.ImportFormRequest) error {
	return receiver.updateFormSettingByType(req, value.SettingTypeImportSignUpForms)
}

func (receiver *SettingRepository) GetRegistrationPreset1Setting() (*entity.SSetting, error) {
	return receiver.getSettingsByType(value.SettingTypeSignUpPresetValue1)
}

func (receiver *SettingRepository) UpdateRegistrationPreset1(url string) (*entity.SSetting, error) {
	setting, err := receiver.getSettingsByType(value.SettingTypeSignUpPresetValue1)
	if setting == nil {
		importSetting := SignUpFormSetting{
			SpreadSheetId: url,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		result := receiver.DBConn.Create(&entity.SSetting{
			Settings: datatypes.JSON(string(b)),
			Type:     value.SettingTypeSignUpPresetValue1,
		})

		return nil, result.Error
	} else {
		importSetting := SignUpFormSetting{
			SpreadSheetId: url,
		}
		b, err := json.Marshal(importSetting)
		if err != nil {
			return nil, err
		}
		setting.Settings = datatypes.JSON(string(b))

		if err = receiver.DBConn.Save(&setting).Error; err != nil {
			return nil, err
		}
	}

	return setting, err
}

func FindDeviceSyncSetting(conn *gorm.DB) (entity.SSetting, error) {
	var setting entity.SSetting
	err := conn.Where("type = ?", value.SettingTypeSyncDevices).First(&setting).Error
	if err != nil {
		return setting, err
	}

	return setting, nil
}
