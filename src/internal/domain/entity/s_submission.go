package entity

import (
	"gorm.io/datatypes"
	"sen-global-api/internal/domain/value"
	"time"
)

type Messaging struct {
	Email        []string `json:"email" binding:"required"`
	Value3       []string `json:"value3" binding:"required"`
	MessageBox   *string  `json:"messageBox"`
	QuestionType string   `json:"questionType" binding:"required"`
}

type SubmissionDataItem struct {
	QuestionId string     `json:"question_id" binding:"required"`
	Question   string     `json:"question" binding:"required"`
	Answer     string     `json:"answer" binding:"required"`
	Messaging  *Messaging `json:"messaging"`
}

type SubmissionData struct {
	Items []SubmissionDataItem `json:"items" binding:"required"`
}

type SSubmission struct {
	ID                 uint64                 `gorm:"primary_key;auto_increment;not null"`
	FormId             uint64                 `gorm:"not null"`
	FormName           string                 `gorm:"type:varchar(255);not null;default:''"`
	FormNote           string                 `gorm:"type:varchar(255);not null;default:''"`
	FormSpreadsheetUrl string                 `gorm:"type:varchar(255);not null;default:''"`
	DeviceId           string                 `gorm:"not null"`
	DeviceFirstValue   string                 `gorm:"type:varchar(255);not null"`
	DeviceSecondValue  string                 `gorm:"type:varchar(255);not null"`
	DeviceThirdValue   string                 `gorm:"type:varchar(255);not null"`
	DeviceName         string                 `gorm:"type:varchar(255);not null;default:''"`
	DeviceNote         string                 `gorm:"type:varchar(255);not null:default:''"`
	SpreadsheetId      string                 `gorm:"type:varchar(255);not null"`
	SheetName          string                 `gorm:"type:varchar(255);not null"`
	SubmissionData     datatypes.JSON         `gorm:"type:json;not null;default:'{}'"`
	SubmissionType     value.SubmissionType   `gorm:"not null"`
	Status             value.SubmissionStatus `gorm:"not null;default:1"`
	OpenedAt           time.Time              `gorm:"default:CURRENT_TIMESTAMP;not null"`
	NumberAttempt      uint64                 `gorm:"not null;default:0"`
	CreatedAt          time.Time              `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt          time.Time              `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
