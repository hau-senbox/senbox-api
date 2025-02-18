package request

type UpdateSettingNameRequest struct {
	FormSetting                     *string `json:"form_setting_1"`
	FormSetting2                    *string `json:"form_setting_2"`
	FormSetting3                    *string `json:"form_setting_3"`
	FormSetting4                    *string `json:"form_setting_4"`
	UrlSetting                      *string `json:"url_setting"`
	OutputSetting                   *string `json:"output_setting"`
	SummarySetting                  *string `json:"summary_setting"`
	SyncDevicesSetting              *string `json:"sync_devices_setting"`
	SyncToDosSetting                *string `json:"sync_to_dos_setting"`
	EmailSetting                    *string `json:"email_setting"`
	OutputTemplateSetting           *string `json:"output_template_setting"`
	OutputTemplateForTeacherSetting *string `json:"output_template_for_teacher_setting"`
	SignUpButton1Setting            *string `json:"sign_up_button_1_setting"`
	SignUpButton2Setting            *string `json:"sign_up_button_2_setting"`
	SignUpButton3Setting            *string `json:"sign_up_button_3_setting"`
	SignUpButton4Setting            *string `json:"sign_up_button_4_setting"`
	SignUpButton5Setting            *string `json:"sign_up_button_5_setting"`
	SignUpButtonConfiguration       *string `json:"sign_up_button_configuration_setting"`
	RegistrationFormSetting         *string `json:"registration_form_setting"`
	RegistrationSubmissionSetting   *string `json:"registration_submission_setting"`
	RegistrationPreset2Setting      *string `json:"registration_preset_2_setting"`
	APIDistributerSetting           *string `json:"api_distributer_setting"`
	CodeCountingDataSetting         *string `json:"code_counting_data_setting"`
	ImportSignUpFormsSetting        *string `json:"import_sign_up_forms_setting"`
	RegistrationPreset1Setting      *string `json:"registration_preset_1_setting"`
}
