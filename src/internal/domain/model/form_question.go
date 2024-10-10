package model

import (
	"gorm.io/datatypes"
	"sen-global-api/internal/domain/value"
	"time"
)

type FormQuestionItem struct {
	QuestionId     string                  `gorm:"type:varchar(255);primary_key;not null"`
	QuestionName   string                  `gorm:"type:varchar(255);not null;default:''"`
	QuestionType   string                  `gorm:"type:varchar(255);not null;default:1"`
	Question       string                  `gorm:"type:varchar(255);not null;default:''"`
	Attributes     datatypes.JSON          `gorm:"type:json;not null;default:'{}'"`
	Status         value.Status            `gorm:"type:int;not null;default:0"`
	Order          int                     `gorm:"type:int;not null;default:0"`
	AnswerRequired bool                    `gorm:"type:tinyint(1);not null;default:0"`
	EnableOnMobile value.QuestionForMobile `gorm:"type:varchar(16);not null;default:'enabled'"`
	QuestionUniqueId *string        `gorm:"type:varchar(255);default:null"`
	CreatedAt      time.Time               `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt      time.Time               `gorm:"default:CURRENT_TIMESTAMP;not null"`
}
