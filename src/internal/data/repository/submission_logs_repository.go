package repository

import (
	"context"
	"sen-global-api/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

type SubmissionLogsRepository struct {
	DBConn *gorm.DB
}

func (r *SubmissionLogsRepository) GetSubmissionsFormLogs(
	ctx context.Context,
	start, end *time.Time,
	qrCode string,
	page, limit int,
) ([]entity.DataLog, int, error) {

	var logs []entity.DataLog
	var total int64

	query := r.DBConn.WithContext(ctx).
		Model(&entity.DataLog{}).
		Where("endpoint = ?", "/v1/form")

	if start != nil && end != nil {
		query = query.Where("created_at BETWEEN ? AND ?", start, end)
	} else if start != nil && end == nil {
		query = query.Where("created_at >= ?", start)
	} else if start == nil && end != nil {
		query = query.Where("created_at <= ?", end)
	}

	if qrCode != "" {
		query = query.Where("JSON_EXTRACT(payload, '$.qr_code') = ?", qrCode)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, int(total), nil
}

func (r *SubmissionLogsRepository) GetSubmissionsFormLogsBySubmit(
	ctx context.Context,
	start, end *time.Time,
	qrCode string,
	customID string,
	page, limit int,
) ([]entity.DataLog, int, error) {

	var logs []entity.DataLog
	var total int64

	query := r.DBConn.WithContext(ctx).
		Model(&entity.DataLog{}).
		Where("endpoint = ?", "/v1/form/submit")

	if start != nil && end != nil {
		query = query.Where("created_at BETWEEN ? AND ?", start, end)
	} else if start != nil && end == nil {
		query = query.Where("created_at >= ?", start)
	} else if start == nil && end != nil {
		query = query.Where("created_at <= ?", end)
	}

	if qrCode != "" {
		query = query.Where("JSON_EXTRACT(payload, '$.qr_code') = ?", qrCode)
	}

	if customID != "" {
		query = query.Where("JSON_EXTRACT(payload, '$.student_custom_id') = ?", customID)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page > 0 && limit > 0 {
		offset := (page - 1) * limit
		query = query.Offset(offset).Limit(limit)
	}

	if err := query.Order("created_at DESC").Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, int(total), nil

}
