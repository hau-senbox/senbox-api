package repository

import (
	"encoding/json"
	"sen-global-api/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

type SubmissionRepository struct {
	DBConn *gorm.DB
}

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
type CreateSubmissionParams struct {
	FormId         uint64
	UserId         string
	SubmissionData SubmissionData
	OpenedAt       time.Time
}

func (receiver *SubmissionRepository) CreateSubmission(params CreateSubmissionParams) error {
	items := make([]entity.SubmissionDataItem, 0)
	for _, item := range params.SubmissionData.Items {
		var msg *entity.Messaging = nil
		if item.Messaging != nil {
			msg = &entity.Messaging{
				Email:        item.Messaging.Email,
				Value3:       item.Messaging.Value3,
				MessageBox:   item.Messaging.MessageBox,
				QuestionType: item.Messaging.QuestionType,
			}
		}
		items = append(items, entity.SubmissionDataItem{
			QuestionId: item.QuestionId,
			Question:   item.Question,
			Answer:     item.Answer,
			Messaging:  msg,
		})
	}

	data := entity.SubmissionData{
		Items: items,
	}

	dataInJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	submission := entity.SSubmission{
		FormId:         params.FormId,
		UserId:         params.UserId,
		SubmissionData: dataInJSON,
		OpenedAt:       params.OpenedAt,
	}

	return receiver.DBConn.Create(&submission).Error
}

func (receiver *SubmissionRepository) FindRecentByFormId(formId uint64, userId string) (entity.SSubmission, error) {
	var submission entity.SSubmission
	err := receiver.DBConn.Where("form_id = ? AND user_id = ?", formId, userId).Order("created_at DESC").First(&submission).Error
	if err != nil {
		return entity.SSubmission{}, err
	}

	return submission, err
}

func (receiver *SubmissionRepository) DublicateSubmissions(params CreateSubmissionParams) error {
	items := make([]entity.SubmissionDataItem, 0)
	for _, item := range params.SubmissionData.Items {
		var msg *entity.Messaging = nil
		if item.Messaging != nil {
			msg = &entity.Messaging{
				Email:        item.Messaging.Email,
				Value3:       item.Messaging.Value3,
				MessageBox:   item.Messaging.MessageBox,
				QuestionType: item.Messaging.QuestionType,
			}
		}
		items = append(items, entity.SubmissionDataItem{
			QuestionId: item.QuestionId,
			Question:   item.Question,
			Answer:     item.Answer,
			Messaging:  msg,
		})
	}

	data := entity.SubmissionData{
		Items: items,
	}

	dataInJSON, err := json.Marshal(data)
	if err != nil {
		return err
	}

	submission := entity.SSubmission{
		FormId:         params.FormId,
		UserId:         params.UserId,
		SubmissionData: dataInJSON,
		OpenedAt:       params.OpenedAt,
	}

	return receiver.DBConn.Create(&submission).Error
}
