package usecase

import (
	"encoding/json"
	"sen-global-api/internal/data/repository"

	log "github.com/sirupsen/logrus"
)

type GetSettingsUseCase struct {
	*repository.SettingRepository
}

type ImportSetting struct {
	SettingName    string `json:"setting_name"`
	SpreadSheetUrl string `json:"spreadsheet_url"`
	AutoImport     bool   `json:"auto"`
	Interval       uint64 `json:"interval"`
}

type OutputSetting struct {
	FolderId    string `json:"folder_id"`
	SettingName string `json:"setting_name"`
}

type SummarySetting struct {
	SpreadsheetId string `json:"spreadsheet_id"`
	SettingName   string `json:"setting_name"`
}

type APIDistributerSetting struct {
	SettingName    string `json:"setting_name"`
	SpreadSheetUrl string `json:"spreadsheet_url"`
}

type SignUpTextButtonSetting struct {
	Name        string `json:"name"`
	Value       string `json:"value"`
	SettingName string `json:"setting_name"`
}

type AppSettings struct {
	Form                      *ImportSetting
	Form2                     *ImportSetting
	Form3                     *ImportSetting
	Form4                     *ImportSetting
	Url                       *ImportSetting
	Output                    *OutputSetting
	Summary                   *SummarySetting
	SyncDevices               *ImportSetting
	SyncToDos                 *ImportSetting
	EmailSetting              *SummarySetting
	OutputTemplate            *SummarySetting
	OutputTemplateForTeacher  *SummarySetting
	SignUpButton1             *SignUpTextButtonSetting
	SignUpButton2             *SignUpTextButtonSetting
	SignUpButton3             *SignUpTextButtonSetting
	SignUpButton4             *SignUpTextButtonSetting
	SignUpButton5             *SignUpTextButtonSetting
	SignUpButtonConfiguration *SummarySetting
	RegistrationForm          *SummarySetting
	RegistrationSubmission    *SummarySetting
	RegistrationPreset2       *SummarySetting
	APIDistributer            *APIDistributerSetting
	CodeCountingData          *APIDistributerSetting
	SignUpForms               *ImportSetting
	RegistrationPreset1       *SummarySetting
}

func (receiver *GetSettingsUseCase) GetSettings() (*AppSettings, error) {
	formSettingsData, err := receiver.GetFormSettings()
	if err != nil {
		log.Info(err.Error())
	}

	formSettingsData2, err := receiver.GetFormSettings2()
	if err != nil {
		log.Info(err.Error())
	}

	formSettingsData3, err := receiver.GetFormSettings3()
	if err != nil {
		log.Info(err.Error())
	}

	formSettingsData4, err := receiver.GetFormSettings4()
	if err != nil {
		log.Info(err.Error())
	}

	urlSettingsData, err := receiver.GetUrlSettings()
	if err != nil {
		log.Info(err.Error())
	}
	outputSettingsData, err := receiver.GetOutputSettings()
	if err != nil {
		log.Info(err.Error())
	}

	summarySettingsData, err := receiver.GetSummarySettings()
	if err != nil {
		log.Info(err.Error())
	}

	syncDevicesSettingsData, err := receiver.GetSyncDevicesSettings()
	if err != nil {
		log.Info(err.Error())
	}

	syncToDosSettingsData, err := receiver.GetSyncToDosSettings()
	if err != nil {
		log.Info(err.Error())
	}

	emailSettingsData, err := receiver.GetEmailSettings()
	if err != nil {
		log.Info(err.Error())
	}

	outputTemplateSettingData, err := receiver.GetOutputTemplateSettings()
	if err != nil {
		log.Info(err.Error())
	}

	outputTemplateForTeacherSettingData, err := receiver.GetOutputTemplateSettingsForTeacher()
	if err != nil {
		log.Info(err.Error())
	}

	signUpButton1SettingData, err := receiver.GetSignUpButton1Setting()
	if err != nil {
		log.Info(err.Error())
	}

	signUpButton2SettingData, err := receiver.GetSignUpButton2Setting()
	if err != nil {
		log.Info(err.Error())
	}

	signUpButton3SettingData, err := receiver.GetSignUpButton3Setting()
	if err != nil {
		log.Info(err.Error())
	}

	signUpButton4SettingData, err := receiver.GetSignUpButton4Setting()
	if err != nil {
		log.Info(err.Error())
	}

	signUpButton5SettingData, err := receiver.GetSignUpButton5Setting()
	if err != nil {
		log.Info(err.Error())
	}

	signUpButtonConfigurationSettingData, err := receiver.GetSignUpButtonConfigurationSetting()
	if err != nil {
		log.Info(err.Error())
	}

	registrationFormSettingData, err := receiver.GetRegistrationFormSetting()
	if err != nil {
		log.Info(err.Error())
	}

	registrationSubmissionSettingData, err := receiver.GetRegistrationSubmissionSetting()
	if err != nil {
		log.Info(err.Error())
	}

	registrationPreset2SettingData, err := receiver.GetRegistrationPreset2Setting()
	if err != nil {
		log.Info(err.Error())
	}

	apiDistributerSettingData, err := receiver.GetAPIDistributerSetting()
	if err != nil {
		log.Info(err.Error())
	}

	codeCountingData, err := receiver.GetCodeCountingDataSetting()
	if err != nil {
		log.Info(err.Error())
	}

	signUpFormSettingData, err := receiver.GetImportSignUpFormsSetting()
	if err != nil {
		log.Info(err.Error())
	}

	registrationPreset1SettingData, err := receiver.GetRegistrationPreset1Setting()
	if err != nil {
		log.Info(err.Error())
	}

	var formSettings *ImportSetting = nil
	var formSettings2 *ImportSetting = nil
	var formSettings3 *ImportSetting = nil
	var formSettings4 *ImportSetting = nil
	var urlSettings *ImportSetting = nil
	var outputSettings *OutputSetting = nil
	var summarySettings *SummarySetting = nil
	var syncDevicesSettings *ImportSetting = nil
	var syncToDosSettings *ImportSetting = nil
	var emailSettings *SummarySetting = nil
	var outputTemplateSettings *SummarySetting = nil
	var outputTemplateForTeacherSettings *SummarySetting = nil
	var signUpButton1Settings *SignUpTextButtonSetting = nil
	var signUpButton2Settings *SignUpTextButtonSetting = nil
	var signUpButton3Settings *SignUpTextButtonSetting = nil
	var signUpButton4Settings *SignUpTextButtonSetting = nil
	var signUpButton5Settings *SignUpTextButtonSetting = nil
	var signUpButtonConfigurationSettings *SummarySetting = nil
	var registrationFormSettings *SummarySetting = nil
	var registrationSubmissionSettings *SummarySetting = nil
	var registrationPreset2Settings *SummarySetting = nil
	var apiDistributerSettings *APIDistributerSetting = nil
	var codeCountingSettings *APIDistributerSetting = nil
	var signUpFormSettings *ImportSetting = nil
	var registrationPreset1Settings *SummarySetting = nil

	if formSettingsData != nil {
		err = json.Unmarshal(formSettingsData.Settings, &formSettings)
		if err != nil {
			log.Info(err.Error())
		}
		formSettings.SettingName = formSettingsData.SettingName
	}

	if formSettingsData2 != nil {
		err = json.Unmarshal(formSettingsData2.Settings, &formSettings2)
		if err != nil {
			log.Info(err.Error())
		}
		formSettings2.SettingName = formSettingsData2.SettingName
	}

	if formSettingsData3 != nil {
		err = json.Unmarshal(formSettingsData3.Settings, &formSettings3)
		if err != nil {
			log.Info(err.Error())
		}
		formSettings3.SettingName = formSettingsData3.SettingName
	}

	if formSettingsData4 != nil {
		err = json.Unmarshal(formSettingsData4.Settings, &formSettings4)
		if err != nil {
			log.Info(err.Error())
		}
		formSettings4.SettingName = formSettingsData4.SettingName
	}

	if urlSettingsData != nil {
		err = json.Unmarshal(urlSettingsData.Settings, &urlSettings)
		if err != nil {
			log.Info(err.Error())
		}
		urlSettings.SettingName = urlSettingsData.SettingName
	}

	if outputSettingsData != nil {
		err = json.Unmarshal(outputSettingsData.Settings, &outputSettings)
		if err != nil {
			log.Info(err.Error())
		}
		outputSettings.SettingName = outputSettingsData.SettingName
	}

	if summarySettingsData != nil {
		err = json.Unmarshal(summarySettingsData.Settings, &summarySettings)
		if err != nil {
			log.Info(err.Error())
		}
		summarySettings.SettingName = summarySettingsData.SettingName
	}

	if syncDevicesSettingsData != nil {
		err = json.Unmarshal(syncDevicesSettingsData.Settings, &syncDevicesSettings)
		if err != nil {
			log.Info(err.Error())
		}
		syncDevicesSettings.SettingName = syncDevicesSettingsData.SettingName
	}

	if syncToDosSettingsData != nil {
		err = json.Unmarshal(syncToDosSettingsData.Settings, &syncToDosSettings)
		if err != nil {
			log.Info(err.Error())
		}
		syncToDosSettings.SettingName = syncToDosSettingsData.SettingName
	}

	if emailSettingsData != nil {
		err = json.Unmarshal(emailSettingsData.Settings, &emailSettings)
		if err != nil {
			log.Info(err.Error())
		}
		emailSettings.SettingName = emailSettingsData.SettingName
	}

	if outputTemplateSettingData != nil {
		err = json.Unmarshal(outputTemplateSettingData.Settings, &outputTemplateSettings)
		if err != nil {
			log.Info(err.Error())
		}
		outputTemplateSettings.SettingName = outputTemplateSettingData.SettingName
	}

	if outputTemplateForTeacherSettingData != nil {
		err = json.Unmarshal(outputTemplateForTeacherSettingData.Settings, &outputTemplateForTeacherSettings)
		if err != nil {
			log.Info(err.Error())
		}
		outputTemplateForTeacherSettings.SettingName = outputTemplateForTeacherSettingData.SettingName
	}

	if signUpButton1SettingData != nil {
		err = json.Unmarshal(signUpButton1SettingData.Settings, &signUpButton1Settings)
		if err != nil {
			log.Info(err.Error())
		}
		signUpButton1Settings.SettingName = signUpButton1SettingData.SettingName
	}

	if signUpButton2SettingData != nil {
		err = json.Unmarshal(signUpButton2SettingData.Settings, &signUpButton2Settings)
		if err != nil {
			log.Info(err.Error())
		}
		signUpButton2Settings.SettingName = signUpButton2SettingData.SettingName
	}

	if signUpButton3SettingData != nil {
		err = json.Unmarshal(signUpButton3SettingData.Settings, &signUpButton3Settings)
		if err != nil {
			log.Info(err.Error())
		}
		signUpButton3Settings.SettingName = signUpButton3SettingData.SettingName
	}

	if signUpButton4SettingData != nil {
		err = json.Unmarshal(signUpButton4SettingData.Settings, &signUpButton4Settings)
		if err != nil {
			log.Info(err.Error())
		}
		signUpButton4Settings.SettingName = signUpButton4SettingData.SettingName
	}

	if signUpButton5SettingData != nil {
		err = json.Unmarshal(signUpButton5SettingData.Settings, &signUpButton5Settings)
		if err != nil {
			log.Info(err.Error())
		}
		signUpButton5Settings.SettingName = signUpButton5SettingData.SettingName
	}

	if signUpButtonConfigurationSettingData != nil {
		err = json.Unmarshal(signUpButtonConfigurationSettingData.Settings, &signUpButtonConfigurationSettings)
		if err != nil {
			log.Info(err.Error())
		}
		signUpButtonConfigurationSettings.SettingName = signUpButtonConfigurationSettingData.SettingName
	}

	if registrationFormSettingData != nil {
		err = json.Unmarshal(registrationFormSettingData.Settings, &registrationFormSettings)
		if err != nil {
			log.Info(err.Error())
		}
		registrationFormSettings.SettingName = registrationFormSettingData.SettingName
	}

	if registrationSubmissionSettingData != nil {
		err = json.Unmarshal(registrationSubmissionSettingData.Settings, &registrationSubmissionSettings)
		if err != nil {
			log.Info(err.Error())
		}
		registrationSubmissionSettings.SettingName = registrationSubmissionSettingData.SettingName
	}

	if registrationPreset2SettingData != nil {
		err = json.Unmarshal(registrationPreset2SettingData.Settings, &registrationPreset2Settings)
		if err != nil {
			log.Info(err.Error())
		}
		registrationPreset2Settings.SettingName = registrationPreset2SettingData.SettingName
	}

	if apiDistributerSettingData != nil {
		err = json.Unmarshal(apiDistributerSettingData.Settings, &apiDistributerSettings)
		if err != nil {
			log.Info(err.Error())
		}
		apiDistributerSettings.SettingName = apiDistributerSettingData.SettingName
	}

	if codeCountingData != nil {
		err = json.Unmarshal(codeCountingData.Settings, &codeCountingSettings)
		if err != nil {
			log.Info(err.Error())
		}
		codeCountingSettings.SettingName = codeCountingData.SettingName
	}

	if signUpFormSettingData != nil {
		err = json.Unmarshal(signUpFormSettingData.Settings, &signUpFormSettings)
		if err != nil {
			log.Info(err.Error())
		}
		signUpFormSettings.SettingName = signUpFormSettingData.SettingName
	}

	if registrationPreset1SettingData != nil {
		err = json.Unmarshal(registrationPreset1SettingData.Settings, &registrationPreset1Settings)
		if err != nil {
			log.Info(err.Error())
		}
		registrationPreset1Settings.SettingName = registrationPreset1SettingData.SettingName
	}

	return &AppSettings{
		Form:                      formSettings,
		Form2:                     formSettings2,
		Form3:                     formSettings3,
		Form4:                     formSettings4,
		Url:                       urlSettings,
		Output:                    outputSettings,
		Summary:                   summarySettings,
		SyncDevices:               syncDevicesSettings,
		SyncToDos:                 syncToDosSettings,
		EmailSetting:              emailSettings,
		OutputTemplate:            outputTemplateSettings,
		OutputTemplateForTeacher:  outputTemplateForTeacherSettings,
		SignUpButton1:             signUpButton1Settings,
		SignUpButton2:             signUpButton2Settings,
		SignUpButton3:             signUpButton3Settings,
		SignUpButton4:             signUpButton4Settings,
		SignUpButton5:             signUpButton5Settings,
		SignUpButtonConfiguration: signUpButtonConfigurationSettings,
		RegistrationForm:          registrationFormSettings,
		RegistrationSubmission:    registrationSubmissionSettings,
		RegistrationPreset2:       registrationPreset2Settings,
		APIDistributer:            apiDistributerSettings,
		CodeCountingData:          codeCountingSettings,
		SignUpForms:               signUpFormSettings,
		RegistrationPreset1:       registrationPreset1Settings,
	}, nil
}
