package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type CreateUserQuestionParams struct {
	DeviceID   string
	QuestionID string
	Order      int
}

type DeviceQuestionRepository struct {
	DBConn *gorm.DB
}

type UserQuestion struct {
	QuestionID     string `sql:"question_id"`
	Question       string `sql:"question"`
	QuestionType   string `sql:"question_type"`
	Attributes     string `sql:"attributes"`
	Order          int    `sql:"order"`
	EnableOnMobile string `sql:"enable_on_mobile"`
}

func (receiver *DeviceQuestionRepository) GetQuestionsBelongToDevice(device *entity.SDevice) ([]UserQuestion, error) {
	var rawUserQuestions []UserQuestion
	tx := receiver.DBConn.Table("s_device_question").Joins("INNER JOIN s_question ON s_device_question.question_id = s_question.question_id").Where("s_device_question.device_id = ?", device.ID).Select("*").Find(&rawUserQuestions)
	err := tx.Find(&rawUserQuestions).Error
	if err != nil {
		return nil, err
	}

	return rawUserQuestions, nil
}
