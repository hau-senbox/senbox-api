package response

type GetSettingsResponseDataImport struct {
	SettingName    string `json:"setting_name" binding:"required"`
	SpreadSheetUrl string `json:"spread_sheet_url" binding:"required"`
	Auto           bool   `json:"auto" binding:"required"`
	Interval       uint64 `json:"interval_in_minutes" binding:"required"`
}

type GetSettingsResponseDataSubmission struct {
	SettingName string `json:"setting_name" binding:"required"`
	FolderUrl   string `json:"folder_url" binding:"required"`
}

type GetSettingsResponseAPIDistributor struct {
	SettingName    string `json:"setting_name" binding:"required"`
	SpreadSheetUrl string `json:"spread_sheet_url" binding:"required"`
}

type GetSettingsResponseData struct {
	ImportFormsSetting        *GetSettingsResponseDataImport     `json:"import_forms_setting"`
	ImportFormsSetting2       *GetSettingsResponseDataImport     `json:"import_forms_setting_2"`
	ImportFormsSetting3       *GetSettingsResponseDataImport     `json:"import_forms_setting_3"`
	ImportFormsSetting4       *GetSettingsResponseDataImport     `json:"import_forms_setting_4"`
	ImportRedirectUrlsSetting *GetSettingsResponseDataImport     `json:"import_redirect_urls_setting"`
	Output                    *GetSettingsResponseDataSubmission `json:"output"`
	Summary                   *GetSettingsResponseDataSummary    `json:"summary"`
	SyncDevices               *GetSettingsResponseDataImport     `json:"sync_devices"`
	ImportToDoListSetting     *GetSettingsResponseDataImport     `json:"import_todo_list_setting"`
	EmailHistory              *GetSettingsResponseDataSummary    `json:"email_history"`
	OutputTemplate            *GetSettingsResponseDataSummary    `json:"output_template"`
	OutputTemplateForTeacher  *GetSettingsResponseDataSummary    `json:"output_template_for_teacher"`
	SignUpButton1             *GetSettingsResponseTextButton     `json:"sign_up_button_1"`
	SignUpButton2             *GetSettingsResponseTextButton     `json:"sign_up_button_2"`
	SignUpButton3             *GetSettingsResponseTextButton     `json:"sign_up_button_3"`
	SignUpButton4             *GetSettingsResponseTextButton     `json:"sign_up_button_4"`
	SignUpButton5             *GetSettingsResponseTextButton     `json:"sign_up_button_5"`
	SignUpButtonConfiguration *GetSettingsResponseDataSummary    `json:"sign_up_button_configuration"`
	RegistrationForm          *GetSettingsResponseDataSummary    `json:"registration_form"`
	RegistrationSubmission    *GetSettingsResponseDataSummary    `json:"registration_submission"`
	RegistrationPreset2       *GetSettingsResponseDataSummary    `json:"registration_preset_2"`
	APIDistributer            *GetSettingsResponseAPIDistributor `json:"api_distributer"`
	CodeCountingData          *GetSettingsResponseAPIDistributor `json:"code_counting_data"`
	SignUpFormsSetting        *GetSettingsResponseDataImport     `json:"import_signup_forms"`
	RegistrationPreset1       *GetSettingsResponseDataSummary    `json:"registration_preset_1"`
}

type GetSettingsResponse struct {
	Data GetSettingsResponseData `json:"data" binding:"required"`
}

type GetSettingsResponseDataSummary struct {
	SettingName    string `json:"setting_name" binding:"required"`
	SpreadSheetUrl string `json:"spreadsheet_url" binding:"required"`
}

type GetSettingsResponseTextButton struct {
	SettingName string `json:"setting_name" binding:"required"`
	Name        string `json:"name" binding:"required"`
	Value       string `json:"value" binding:"required"`
}
