package entity

import (
	"sen-global-api/internal/domain/value"
	"strings"
	"time"

	"github.com/google/uuid"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type SQuestion struct {
	ID               uuid.UUID               `gorm:"type:char(36);primary_key"`
	Question         string                  `gorm:"type:varchar(1000);not null;default:''"`
	QuestionType     string                  `gorm:"type:varchar(255);not null;default:1"`
	Attributes       datatypes.JSON          `gorm:"type:json;not null;default:'{}'"`
	Status           value.Status            `gorm:"type:int;not null;default:0"`
	Set              string                  `gorm:"type:varchar(255);not null;default:''"`
	EnableOnMobile   value.QuestionForMobile `gorm:"type:varchar(16);not null;default:'enabled'"`
	QuestionUniqueID *string                 `gorm:"type:varchar(255);default:null"`
	Key              string                  `gorm:"type:varchar(255);not null;default:''"`
	DB               string                  `gorm:"type:varchar(255);not null;default:''"`
	CreatedAt        time.Time               `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt        time.Time               `gorm:"default:CURRENT_TIMESTAMP;not null"`
}

func (question *SQuestion) BeforeSave(tx *gorm.DB) (err error) {
	question.QuestionType = strings.ToLower(question.QuestionType)

	return err
}
