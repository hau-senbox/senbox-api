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

type SubmissionDataItem struct {
	QuestionID  string `json:"question_id" binding:"required"`
	QuestionKey string `json:"question_key"`
	QuestionDB  string `json:"question_db"`
	Question    string `json:"question" binding:"required"`
	Answer      string `json:"answer" binding:"required"`
}

type SubmissionData struct {
	Items []SubmissionDataItem `json:"items" binding:"required"`
}
type CreateSubmissionParams struct {
	FormID         uint64
	UserID         string
	SubmissionData SubmissionData
	OpenedAt       time.Time
}
type GetSubmissionByConditionParam struct {
	FormID      uint64
	UserID      string
	QuestionKey string
	QuestionDB  string
}

func (receiver *SubmissionRepository) CreateSubmission(params CreateSubmissionParams) error {
	items := make([]entity.SubmissionDataItem, 0)
	for _, item := range params.SubmissionData.Items {
		items = append(items, entity.SubmissionDataItem{
			QuestionID:  item.QuestionID,
			QuestionKey: item.QuestionKey,
			QuestionDB:  item.QuestionDB,
			Question:    item.Question,
			Answer:      item.Answer,
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
		FormID:         params.FormID,
		UserID:         params.UserID,
		SubmissionData: dataInJSON,
		OpenedAt:       params.OpenedAt,
	}

	return receiver.DBConn.Create(&submission).Error
}

func (receiver *SubmissionRepository) FindRecentByFormID(formID uint64, userID string) (entity.SSubmission, error) {
	var submission entity.SSubmission
	err := receiver.DBConn.Where("form_id = ? AND user_id = ?", formID, userID).Order("created_at DESC").First(&submission).Error
	if err != nil {
		return entity.SSubmission{}, err
	}

	return submission, err
}

func (receiver *SubmissionRepository) DuplicateSubmissions(params CreateSubmissionParams) error {
	items := make([]entity.SubmissionDataItem, 0)
	for _, item := range params.SubmissionData.Items {
		items = append(items, entity.SubmissionDataItem{
			QuestionID: item.QuestionID,
			Question:   item.Question,
			Answer:     item.Answer,
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
		FormID:         params.FormID,
		UserID:         params.UserID,
		SubmissionData: dataInJSON,
		OpenedAt:       params.OpenedAt,
	}

	return receiver.DBConn.Create(&submission).Error
}

func (receiver *SubmissionRepository) GetSubmissionByCondition(param GetSubmissionByConditionParam) ([]SubmissionDataItem, error) {
	var submissions []entity.SSubmission

	// Bước 1: Truy vấn theo form_id và user_id
	err := receiver.DBConn.
		Where("form_id = ? AND user_id = ?", param.FormID, param.UserID).
		Find(&submissions).Error
	if err != nil {
		return nil, err
	}

	var result []SubmissionDataItem

	// Bước 2: Duyệt từng bản ghi submission
	for _, submission := range submissions {
		var data SubmissionData

		if err := json.Unmarshal(submission.SubmissionData, &data); err != nil {
			continue // skip nếu có lỗi parse JSON
		}

		// Bước 3: Lọc dữ liệu theo question_key, question_db nếu được truyền
		for _, item := range data.Items {
			if (param.QuestionKey == "" || item.QuestionKey == param.QuestionKey) &&
				(param.QuestionDB == "" || item.QuestionDB == param.QuestionDB) {
				result = append(result, item)
			}
		}
	}

	return result, nil
}
