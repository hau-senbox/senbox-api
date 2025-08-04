package repository

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type SyncQueueRepository struct {
	DBConn *gorm.DB
}

func (r *SyncQueueRepository) Create(q *entity.SyncQueue) error {
	return r.DBConn.Create(q).Error
}

func (r *SyncQueueRepository) UpdateStatus(id uint64, status string) error {
	return r.DBConn.Model(&entity.SyncQueue{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *SyncQueueRepository) HasPendingQueue() (bool, error) {
	var count int64
	err := r.DBConn.Model(&entity.SyncQueue{}).
		Where("status = ?", value.SyncQueueStatusPending).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count == 0, nil // false nếu có queue đang pending
}

func (r *SyncQueueRepository) GetAll() ([]entity.SyncQueue, error) {
	var queues []entity.SyncQueue
	err := r.DBConn.Find(&queues).Error
	if err != nil {
		return nil, err
	}
	return queues, nil
}

func (r *SyncQueueRepository) UpdateOrCreateBySpreadsheetIDAndSheetName(q *entity.SyncQueue) error {
	return r.DBConn.
		Clauses(clause.OnConflict{
			Columns: []clause.Column{
				{Name: "spreadsheet_id"},
				{Name: "sheet_name"},
			},
			DoUpdates: clause.AssignmentColumns([]string{
				"last_submission_id",
				"last_submitted_at",
				"form_notes",
				"status",
				"updated_at",
			}),
		}).
		Create(q).Error
}
