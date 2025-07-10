package entity

import (
	"time"

	"gorm.io/datatypes"
)

type SubmissionDataItem struct {
	QuestionID string `json:"question_id" binding:"required"`
	Key        string `json:"key"`
	DB         string `json:"db"`
	Question   string `json:"question" binding:"required"`
	Answer     string `json:"answer" binding:"required"`
}

type SubmissionData struct {
	Items []SubmissionDataItem `json:"items" binding:"required"`
}

type SSubmission struct {
	ID             uint64         `gorm:"primary_key;auto_increment;"`
	FormID         uint64         `gorm:"column:form_id;"`
	Form           SForm          `gorm:"foreignKey:FormID;references:id;constraint:OnDelete:CASCADE"`
	UserID         string         `gorm:"column:user_id;"`
	User           SUserEntity    `gorm:"foreignKey:UserID;references:id;constraint:OnDelete:CASCADE"`
	SubmissionData datatypes.JSON `gorm:"type:json;not null;default:'{}'"`
	OpenedAt       time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
	CreatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt      time.Time      `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
