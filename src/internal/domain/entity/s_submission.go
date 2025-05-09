package entity

import (
	"time"

	"gorm.io/datatypes"
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
	ID             uint64         `gorm:"primary_key;auto_increment;"`
	FormId         uint64         `gorm:"column:form_id;"`
	Form           SForm          `gorm:"foreignKey:FormId;references:id;constraint:OnDelete:CASCADE"`
	UserId         string         `gorm:"column:user_id;"`
	User           SUserEntity    `gorm:"foreignKey:UserId;references:id;constraint:OnDelete:CASCADE"`
	SubmissionData datatypes.JSON `gorm:"type:json;not null;default:'{}'"`
	OpenedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
