package entity

import (
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"sen-global-api/internal/domain/value"
	"strings"
	"time"
)

type SQuestion struct {
	QuestionId     string                  `gorm:"type:varchar(255);primary_key;not null"`
	QuestionName   string                  `gorm:"type:varchar(1000);not null;default:''"`
	QuestionType   string                  `gorm:"type:varchar(255);not null;default:1"`
	Question       string                  `gorm:"type:varchar(1000);not null;default:''"`
	Attributes     datatypes.JSON          `gorm:"type:json;not null;default:'{}'"`
	Status         value.Status            `gorm:"type:int;not null;default:0"`
	Set            string                  `gorm:"type:varchar(255);not null;default:''"`
	EnableOnMobile value.QuestionForMobile `gorm:"type:varchar(16);not null;default:'enabled'"`
	QuestionUniqueId *string        `json:"type:varchar(255);default:null"`
	CreatedAt      time.Time               `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt      time.Time               `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func (question *SQuestion) BeforeSave(tx *gorm.DB) (err error) {
	question.QuestionType = strings.ToLower(question.QuestionType)

	return err
}
