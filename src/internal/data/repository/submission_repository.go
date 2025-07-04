package repository

import (
	"encoding/json"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"
	"sort"
	"time"

	"gorm.io/gorm"
)

type SubmissionRepository struct {
	DBConn *gorm.DB
}

type SubmissionDataItem struct {
	SubmissionID string    `json:"id"`
	QuestionID   string    `json:"question_id" binding:"required"`
	QuestionKey  string    `json:"question_key"`
	QuestionDB   string    `json:"question_db"`
	Question     string    `json:"question" binding:"required"`
	Answer       string    `json:"answer" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
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
	TimeSort    value.TimeSort
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

func (receiver *SubmissionRepository) GetSubmissionByCondition(param GetSubmissionByConditionParam) (*SubmissionDataItem, error) {
	var submissions []entity.SSubmission

	query := receiver.DBConn.Where("user_id = ?", param.UserID)

	if param.FormID != 0 {
		query = query.Where("form_id = ?", param.FormID)
	}

	switch param.TimeSort {
	case value.TimeShortOldest:
		query = query.Order("created_at ASC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Find(&submissions).Error
	if err != nil {
		return nil, err
	}

	var result []SubmissionDataItem

	for _, submission := range submissions {
		var data SubmissionData
		if err := json.Unmarshal(submission.SubmissionData, &data); err != nil {
			continue
		}

		for _, item := range data.Items {
			if (item.QuestionKey == param.QuestionKey) &&
				(item.QuestionDB == param.QuestionDB) {
				item.CreatedAt = submission.CreatedAt
				result = append(result, item)
			}
		}
	}

	switch param.TimeSort {
	case value.TimeShortOldest:
		sort.Slice(result, func(i, j int) bool {
			return result[i].CreatedAt.Before(result[j].CreatedAt)
		})
	default:
		sort.Slice(result, func(i, j int) bool {
			return result[i].CreatedAt.After(result[j].CreatedAt)
		})
	}

	if len(result) == 0 {
		return nil, nil
	}

	return &result[0], nil
}
