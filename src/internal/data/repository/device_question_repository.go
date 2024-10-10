package repository

import (
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sen-global-api/internal/domain/entity"
)

type CreateUserQuestionParams struct {
	DeviceId   string
	QuestionId string
	Order      int
}

type DeviceQuestionRepository struct {
	DBConn *gorm.DB
}

type UserQuestion struct {
	QuestionId     string `sql:"question_id"`
	QuestionName   string `sql:"question_name"`
	QuestionType   string `sql:"question_type"`
	Question       string `sql:"question"`
	Attributes     string `sql:"attributes"`
	Order          int    `sql:"order"`
	EnableOnMobile string `sql:"enable_on_mobile"`
}

func (receiver *DeviceQuestionRepository) GetQuestionsBelongToDevice(user *entity.SDevice) ([]UserQuestion, error) {
	var rawUserQuestions []UserQuestion
	tx := receiver.DBConn.Table("s_device_question").Joins("INNER JOIN s_question ON s_device_question.question_id = s_question.question_id").Where("s_device_question.device_id = ?", user.DeviceId).Select("*").Find(&rawUserQuestions)
	err := tx.Find(&rawUserQuestions).Error
	if err != nil {
		return nil, err
	}
	log.Debug("rawUserQuestions: ", rawUserQuestions)
	return rawUserQuestions, nil
}
