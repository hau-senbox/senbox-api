package entity

import (
	"sen-global-api/internal/domain/value"
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SSetting struct {
	ID           int               `gorm:"column:id;primaryKey;autoIncrement"`
	SettingName  string            `gorm:"type:text;not null;default:''"`
	Settings     datatypes.JSON    `gorm:"column:settings;type:json;not null;default:'{}'"`
	Type         value.SettingType `gorm:"type:int;not null;default:0;unique"`
	IntegerValue uint64            `gorm:"column:integer_value;type:int;not null;default:0"`
	CreatedAt    time.Time         `gorm:"column:created_at;default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt    time.Time         `gorm:"column:updated_at;default:CURRENT_TIMESTAMP;not null"`
}

func (s *SSetting) BeforeSave(tx *gorm.DB) (err error) {
	if s.SettingName != "" {
		return
	}
	switch s.Type {
	case value.SettingTypeImportForms:
		s.SettingName = "Form Uploader 1"
	case value.SettingTypeImportForms2:
		s.SettingName = "Form Uploader 2"
	case value.SettingTypeImportForms3:
		s.SettingName = "Form Uploader 3"
	case value.SettingTypeImportForms4:
		s.SettingName = "Form Uploader 4"
	case value.SettingTypeEmailHistory:
		s.SettingName = "Email History"
	case value.SettingTypeImportUrls:
		s.SettingName = "Redirect URLs"
	case value.SettingTypeOutputTemplate:
		s.SettingName = "Output Template"
	case value.SettingTypeOutputTemplateTeacher:
		s.SettingName = "Output Template Teacher"
	case value.SettingTypeSignUpButton1:
		s.SettingName = "Sign Up Button 1"
	case value.SettingTypeSignUpButton2:
		s.SettingName = "Sign Up Button 2"
	case value.SettingTypeSignUpButton3:
		s.SettingName = "Sign Up Button 3"
	case value.SettingTypeSignUpButton4:
		s.SettingName = "Sign Up Button 4"
	case value.SettingTypeSignUpForm:
		s.SettingName = "Sign Up Form"
	case value.SettingTypeSignUpOutput:
		s.SettingName = "Sign Up Output"
	case value.SettingTypeSignUpPresetValue2:
		s.SettingName = "Sign Up Preset Value12"
	case value.SettingTypeSubmission:
		s.SettingName = "Submission"
	case value.SettingTypeSummary:
		s.SettingName = "Summary"
	case value.SettingTypeSyncDevices:
		s.SettingName = "Sync Devices"
	case value.SettingTypeSyncToDos:
		s.SettingName = "Sync ToDos"
	case value.SettingTypeAPIDistributer:
		s.SettingName = "API Distribution"
	case value.SettingTypeCodeCountingData:
		s.SettingName = "Counting Data"
	case value.SettingTypeLogoRefreshInterval:
		s.SettingName = "Logo Refresh Interval(in seconds)"
	case value.SettingTypeImportSignUpForms:
		s.SettingName = "Import Sign Up Forms"
	}
	return
}
