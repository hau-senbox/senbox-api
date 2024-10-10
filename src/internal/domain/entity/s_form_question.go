package entity

type SFormQuestion struct {
	FormId         uint64    `gorm:"not null;primary_key"`
	QuestionId     string    `gorm:"type:varchar(255);not null;primary_key"`
	CreatedAt      string    `gorm:"default:CURRENT_TIMESTAMP;not null"`
	UpdatedAt      string    `gorm:"default:CURRENT_TIMESTAMP;not null"`
	Order          int       `gorm:"type:int;not null;default:0"`
	AnswerRequired bool      `gorm:"type:tinyint(1);not null;default:0"`
	Form           SForm     `gorm:"constraint:OnDelete:CASCADE;"`
	Question       SQuestion `gorm:"constraint:OnDelete:CASCADE;"`
}
