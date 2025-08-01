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
	Key          string    `json:"key"`
	DB           string    `json:"db"`
	Question     string    `json:"question" binding:"required"`
	Answer       string    `json:"answer" binding:"required"`
	CreatedAt    time.Time `json:"created_at"`
}

type SubmissionData struct {
	Items []SubmissionDataItem `json:"items" binding:"required"`
}

type CreateSubmissionParams struct {
	FormID          uint64
	UserID          string
	SubmissionData  SubmissionData
	OpenedAt        time.Time
	StudentCustomID *string
	UserCustomID    *string
}

type GetSubmissionByConditionParam struct {
	FormID       uint64
	UserID       string
	Key          *string
	DB           *string
	TimeSort     value.TimeSort
	DateDuration *value.TimeRange
	Quantity     *string
}

type GetSubmission4MemoriesFormParam struct {
	FormID uint64
	UserId string
}

func (receiver *SubmissionRepository) CreateSubmission(params CreateSubmissionParams) (uint64, error) {
	items := make([]entity.SubmissionDataItem, 0)
	for _, item := range params.SubmissionData.Items {
		items = append(items, entity.SubmissionDataItem{
			QuestionID: item.QuestionID,
			Key:        item.Key,
			DB:         item.DB,
			Question:   item.Question,
			Answer:     item.Answer,
		})
	}

	data := entity.SubmissionData{
		Items: items,
	}

	dataInJSON, err := json.Marshal(data)
	if err != nil {
		return 0, err
	}

	var studentCustomID string
	var userCustomID string
	if params.StudentCustomID != nil {
		studentCustomID = *params.StudentCustomID
	} else {
		studentCustomID = ""
	}
	if params.UserCustomID != nil {
		userCustomID = *params.UserCustomID
	} else {
		userCustomID = ""
	}

	submission := entity.SSubmission{
		FormID:          params.FormID,
		UserID:          params.UserID,
		SubmissionData:  dataInJSON,
		OpenedAt:        params.OpenedAt,
		StudentCustomID: studentCustomID,
		UserCustomID:    userCustomID,
	}

	if err := receiver.DBConn.Create(&submission).Error; err != nil {
		return 0, err
	}

	return submission.ID, nil
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

	if param.DateDuration != nil {
		query = query.Where("created_at BETWEEN ? AND ?", param.DateDuration.Start, param.DateDuration.End)
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
			// Ưu tiên lọc theo Key nếu có
			if item.Key == *param.Key && item.DB == *param.DB {
				if !seen[item.SubmissionID] {
					item.CreatedAt = submission.CreatedAt
					result = append(result, item)
					seen[item.SubmissionID] = true
				}
			}

			// Ưu tiên lọc theo Key nếu có
			// if param.Key != nil && item.Key == *param.Key {
			// 	if !seen[item.SubmissionID] {
			// 		item.CreatedAt = submission.CreatedAt
			// 		result = append(result, item)
			// 		seen[item.SubmissionID] = true
			// 	}
			// }

			// // Nếu có truyền thêm DB, vẫn lọc, nhưng không thêm trùng
			// if param.DB != nil && item.DB == *param.DB {
			// 	if !seen[item.SubmissionID] {
			// 		item.CreatedAt = submission.CreatedAt
			// 		result = append(result, item)
			// 		seen[item.SubmissionID] = true
			// 	}
			// }
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

	//neu khong co quantiy retrun gia tri dau tien
	if param.Quantity == nil {
		first := result[:1]
		return &first, nil
	}

	// Giới hạn số lượng kết quả trả về theo Quantity
	limit := 1
	if param.Quantity != nil {
		if *param.Quantity == "all" {
			return &result, nil
		}
		if qty, err := strconv.Atoi(*param.Quantity); err == nil {
			limit = qty
		}
	}

	if limit > len(result) {
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

	if param.DateDuration != nil {
		query = query.Where("created_at BETWEEN ? AND ?", param.DateDuration.Start, param.DateDuration.End)
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
			if param.Key != nil && item.Key == *param.Key {
				matched = true
			} else if param.DB != nil && item.DB == *param.DB {
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

func (r *SubmissionRepository) GetByUserIdAndFormId(userID string, formID uint64) (*entity.SSubmission, error) {
	var submission entity.SSubmission

	err := r.DBConn.Where("user_id = ? AND form_id = ?", userID, formID).
		Order("created_at DESC").
		First(&submission).Error

	if err != nil {
		return nil, err
	}

	return &submission, nil
}

func (r *SubmissionRepository) GetSubmissionByCreatedAtAndForms(createdAfter time.Time, formNotes []string) ([]*entity.SSubmission, error) {
	var submissions []*entity.SSubmission

	// Join bảng s_form và lọc theo created_at
	adjustedCreatedAfter := createdAfter.Add(500 * time.Millisecond)
	db := r.DBConn.
		Joins("JOIN s_form ON s_form.id = s_submission.form_id").
		Preload("Form").
		Where("s_submission.created_at > ?", adjustedCreatedAfter)

	// Lọc theo danh sách formNotes (nếu có)
	if len(formNotes) > 0 {
		db = db.Where("s_form.note IN ?", formNotes)
	}

	err := db.
		Order("s_submission.created_at ASC").
		Find(&submissions).Error

	if err != nil {
		return nil, err
	}

	return submissions, nil
}
