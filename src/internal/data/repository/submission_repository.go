package repository

import (
	"encoding/json"
	"fmt"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"sort"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type SubmissionRepository struct {
	DBConn *gorm.DB
}

type SubmissionDataItem struct {
	SubmissionID uint64    `json:"id"`
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
	QuestionKey *string
	QuestionDB  *string
	TimeSort    value.TimeSort
	Duration    *value.TimeRange
	Quantity    int
}

type GetSubmission4MemoriesFormParam struct {
	FormID uint64
	UserId string
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

func (receiver *SubmissionRepository) GetSubmissionByCondition(param GetSubmissionByConditionParam) (*[]SubmissionDataItem, error) {
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

		seen := make(map[uint64]bool)

		for _, item := range data.Items {
			item.SubmissionID = submission.ID
			// Ưu tiên lọc theo QuestionKey nếu có
			if param.QuestionKey != nil && item.QuestionKey == *param.QuestionKey {
				if !seen[item.SubmissionID] {
					item.CreatedAt = submission.CreatedAt
					result = append(result, item)
					seen[item.SubmissionID] = true
				}
			}

			// Nếu có truyền thêm QuestionDB, vẫn lọc, nhưng không thêm trùng
			if param.QuestionDB != nil && item.QuestionDB == *param.QuestionDB {
				if !seen[item.SubmissionID] {
					item.CreatedAt = submission.CreatedAt
					result = append(result, item)
					seen[item.SubmissionID] = true
				}
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

	//neu khong co quantiy retrun all
	if param.Quantity == 0 {
		return &result, nil
	}
	// Giới hạn số lượng kết quả trả về theo Quantity
	limit := param.Quantity
	if limit <= 0 || limit > len(result) {
		limit = len(result)
	}

	trimmed := result[:limit]
	return &trimmed, nil

}

func (receiver *SubmissionRepository) GetTotalNrSubmissionByCondition(param GetSubmissionByConditionParam) (*response.GetSubmissionTotalNrResponse, error) {
	var submissions []entity.SSubmission

	query := receiver.DBConn.Where("user_id = ?", param.UserID)
	if param.FormID != 0 {
		query = query.Where("form_id = ?", param.FormID)
	}

	if param.Duration != nil {
		query = query.Where("created_at BETWEEN ? AND ?", param.Duration.Start, param.Duration.End)
	}

	err := query.Find(&submissions).Error
	if err != nil {
		return nil, err
	}

	var total float64 = 0

	for _, submission := range submissions {
		var data SubmissionData
		if err := json.Unmarshal(submission.SubmissionData, &data); err != nil {
			continue
		}

		seen := make(map[uint64]bool)

		for _, item := range data.Items {
			item.SubmissionID = submission.ID

			matched := false
			if param.QuestionKey != nil && item.QuestionKey == *param.QuestionKey {
				matched = true
			} else if param.QuestionDB != nil && item.QuestionDB == *param.QuestionDB {
				matched = true
			}

			if matched && !seen[item.SubmissionID] {
				// Parse answer sang số
				value, err := strconv.ParseFloat(item.Answer, 64)
				if err != nil {
					continue // bỏ qua nếu không phải số
				}

				total += value
				seen[item.SubmissionID] = true
			}
		}
	}

	return &response.GetSubmissionTotalNrResponse{
		Total: strconv.FormatFloat(total, 'f', -1, 64),
	}, nil
}

func (receiver *SubmissionRepository) GetSubmission4MemoriesForm(param GetSubmission4MemoriesFormParam) ([]SubmissionDataItem, error) {
	var submission entity.SSubmission

	err := receiver.DBConn.
		Where("user_id = ?", param.UserId).
		Where("form_id = ?", param.FormID).
		Order("created_at DESC").
		First(&submission).Error

	if err != nil {
		return nil, err
	}

	// Parse JSON
	var data SubmissionData
	if err := json.Unmarshal(submission.SubmissionData, &data); err != nil {
		return nil, fmt.Errorf("failed to parse submission data: %v", err)
	}

	// Gắn SubmissionID và CreatedAt
	for i := range data.Items {
		data.Items[i].SubmissionID = submission.ID
		data.Items[i].CreatedAt = submission.CreatedAt
	}

	// Lấy toàn bộ câu hỏi trong form

	// questions, err := receiver.GetQuestionsByFormID(param.FormID)
	// if err != nil {
	// 	return nil, fmt.Errorf("failed to get questions: %v", err)
	// }

	// // Tạo map các question đã có trong submission (để dễ so sánh)
	// existingMap := make(map[string]bool)
	// for _, item := range data.Items {
	// 	existingMap[item.QuestionID] = true
	// }

	// // Thêm các câu hỏi chưa có trong SubmissionData
	// for _, q := range questions {
	// 	if _, exists := existingMap[q.ID]; !exists {
	// 		data.Items = append(data.Items, SubmissionDataItem{
	// 			SubmissionID: submission.ID,
	// 			QuestionID:   q.ID,
	// 			QuestionKey:  q.QuestionKey,
	// 			QuestionDB:   q.QuestionDB,
	// 			Question:     q.Question,
	// 			Answer:       "",
	// 			CreatedAt:    submission.CreatedAt,
	// 		})
	// 	}
	// }

	// // Tuỳ ý: Sắp xếp theo QuestionKey nếu cần
	// sort.Slice(data.Items, func(i, j int) bool {
	// 	return data.Items[i].QuestionKey < data.Items[j].QuestionKey
	// })

	return data.Items, nil
}
