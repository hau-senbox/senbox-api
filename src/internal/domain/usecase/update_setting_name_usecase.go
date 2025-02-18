package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
)

type UpdateSettingNameUseCase struct {
	db         *gorm.DB
	repository *repository.SettingRepository
}

func NewUpdateSettingNameUseCase(db *gorm.DB) *UpdateSettingNameUseCase {
	return &UpdateSettingNameUseCase{
		db:         db,
		repository: repository.NewSettingRepository(db),
	}
}

func (u *UpdateSettingNameUseCase) Execute(req request.UpdateSettingNameRequest) error {
	if req.FormSetting != nil {
		err := u.repository.SetName(*req.FormSetting, value.SettingTypeImportForms)
		if err != nil {
			return err
		}
	}

	if req.FormSetting2 != nil {
		err := u.repository.SetName(*req.FormSetting2, value.SettingTypeImportForms2)
		if err != nil {
			return err
		}
	}

	if req.FormSetting3 != nil {
		err := u.repository.SetName(*req.FormSetting3, value.SettingTypeImportForms3)
		if err != nil {
			return err
		}
	}

	if req.FormSetting4 != nil {
		err := u.repository.SetName(*req.FormSetting4, value.SettingTypeImportForms4)
		if err != nil {
			return err
		}
	}

	if req.UrlSetting != nil {
		err := u.repository.SetName(*req.UrlSetting, value.SettingTypeImportUrls)
		if err != nil {
			return err
		}
	}

	if req.OutputSetting != nil {
		err := u.repository.SetName(*req.OutputSetting, value.SettingTypeSubmission)
		if err != nil {
			return err
		}
	}

	if req.SummarySetting != nil {
		err := u.repository.SetName(*req.SummarySetting, value.SettingTypeSummary)
		if err != nil {
			return err
		}
	}

	if req.SyncDevicesSetting != nil {
		err := u.repository.SetName(*req.SyncDevicesSetting, value.SettingTypeSyncDevices)
		if err != nil {
			return err
		}
	}

	if req.SyncToDosSetting != nil {
		err := u.repository.SetName(*req.SyncToDosSetting, value.SettingTypeSyncToDos)
		if err != nil {
			return err
		}
	}

	if req.EmailSetting != nil {
		err := u.repository.SetName(*req.EmailSetting, value.SettingTypeEmailHistory)
		if err != nil {
			return err
		}
	}

	if req.OutputTemplateSetting != nil {
		err := u.repository.SetName(*req.OutputTemplateSetting, value.SettingTypeOutputTemplate)
		if err != nil {
			return err
		}
	}

	if req.OutputTemplateForTeacherSetting != nil {
		err := u.repository.SetName(*req.OutputTemplateForTeacherSetting, value.SettingTypeOutputTemplateTeacher)
		if err != nil {
			return err
		}
	}

	if req.SignUpButton1Setting != nil {
		err := u.repository.SetName(*req.SignUpButton1Setting, value.SettingTypeSignUpButton1)
		if err != nil {
			return err
		}
	}

	if req.SignUpButton2Setting != nil {
		err := u.repository.SetName(*req.SignUpButton2Setting, value.SettingTypeSignUpButton2)
		if err != nil {
			return err
		}
	}

	if req.SignUpButton3Setting != nil {
		err := u.repository.SetName(*req.SignUpButton3Setting, value.SettingTypeSignUpButton3)
		if err != nil {
			return err
		}
	}

	if req.SignUpButton4Setting != nil {
		err := u.repository.SetName(*req.SignUpButton4Setting, value.SettingTypeSignUpButton4)
		if err != nil {
			return err
		}
	}

	if req.SignUpButton5Setting != nil {
		err := u.repository.SetName(*req.SignUpButton5Setting, value.SettingTypeSignUpButton5)
		if err != nil {
			return err
		}
	}

	if req.SignUpButtonConfiguration != nil {
		err := u.repository.SetName(*req.SignUpButtonConfiguration, value.SettingTypeSignUpButtonConfiguration)
		if err != nil {
			return err
		}
	}

	if req.RegistrationFormSetting != nil {
		err := u.repository.SetName(*req.RegistrationFormSetting, value.SettingTypeSignUpForm)
		if err != nil {
			return err
		}
	}

	if req.RegistrationSubmissionSetting != nil {
		err := u.repository.SetName(*req.RegistrationSubmissionSetting, value.SettingTypeSignUpOutput)
		if err != nil {
			return err
		}
	}

	if req.RegistrationPreset2Setting != nil {
		err := u.repository.SetName(*req.RegistrationPreset2Setting, value.SettingTypeSignUpPresetValue2)
		if err != nil {
			return err
		}
	}

	if req.APIDistributerSetting != nil {
		err := u.repository.SetName(*req.APIDistributerSetting, value.SettingTypeAPIDistributer)
		if err != nil {
			return err
		}
	}

	if req.CodeCountingDataSetting != nil {
		err := u.repository.SetName(*req.CodeCountingDataSetting, value.SettingTypeCodeCountingData)
		if err != nil {
			return err
		}
	}

	if req.ImportSignUpFormsSetting != nil {
		err := u.repository.SetName(*req.ImportSignUpFormsSetting, value.SettingTypeImportSignUpForms)
		if err != nil {
			return err
		}
	}

	return nil
}
