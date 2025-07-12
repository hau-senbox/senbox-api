package repository

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AnswerRepository struct {
	DBConn *gorm.DB
}

// Create: thêm câu trả lời mới
func (r *AnswerRepository) Create(answer *entity.SAnswer) error {
	return r.DBConn.Create(answer).Error
}

// FindByID: tìm theo ID
func (r *AnswerRepository) FindByID(id uuid.UUID) (*entity.SAnswer, error) {
	var answer entity.SAnswer
	if err := r.DBConn.First(&answer, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &answer, nil
}

// FindBySubmissionID: lấy tất cả câu trả lời theo submission
func (r *AnswerRepository) FindBySubmissionID(submissionID string) ([]entity.SAnswer, error) {
	var answers []entity.SAnswer
	if err := r.DBConn.Where("submission_id = ?", submissionID).Find(&answers).Error; err != nil {
		return nil, err
	}
	return answers, nil
}

// Update: cập nhật câu trả lời
func (r *AnswerRepository) Update(answer *entity.SAnswer) (*entity.SAnswer, error) {
	if err := r.DBConn.Save(answer).Error; err != nil {
		return nil, err
	}
	return answer, nil
}

// Delete: xóa câu trả lời theo ID
func (r *AnswerRepository) Delete(id uuid.UUID) error {
	return r.DBConn.Delete(&entity.SAnswer{}, "id = ?", id).Error
}

// FindByKeyAndDB: lấy danh sách câu trả lời theo Key và DB
func (r *AnswerRepository) FindByKeyAndDB(param GetSubmissionByConditionParam) ([]entity.SAnswer, error) {
	var answers []entity.SAnswer
	query := r.DBConn.Where("`key` = ? AND `db` = ?", param.Key, param.DB)

	if param.DateDuration != nil {
		query = query.Where("created_at BETWEEN ? AND ?", param.DateDuration.Start, param.DateDuration.End)
	}

	switch param.TimeSort {
	case value.TimeShortOldest:
		query = query.Order("created_at ASC")
	default:
		query = query.Order("created_at DESC")
	}

	err := query.Find(&answers).Error
	if err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *AnswerRepository) GetTotalByKeyAndDb(param GetSubmissionByConditionParam) ([]entity.SAnswer, error) {
	var answers []entity.SAnswer
	query := r.DBConn.Where("`key` = ? AND `db` = ?", param.Key, param.DB)

	if param.DateDuration != nil {
		query = query.Where("created_at BETWEEN ? AND ?", param.DateDuration.Start, param.DateDuration.End)
	}

	err := query.Find(&answers).Error
	if err != nil {
		return nil, err
	}
	return answers, nil
}

func (r *AnswerRepository) GetChartTotalByKeyAndDb(param GetSubmissionByConditionParam) ([]entity.SAnswer, error) {
	var answers []entity.SAnswer

	query := r.DBConn.Where("`key` = ? AND `db` = ?", param.Key, param.DB)

	if param.DateDuration != nil {
		query = query.Where("created_at BETWEEN ? AND ?", param.DateDuration.Start, param.DateDuration.End)
	}

	err := query.Order("created_at ASC").Find(&answers).Error

	if err != nil {
		return nil, err
	}

	return answers, nil
}
