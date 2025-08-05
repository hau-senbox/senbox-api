package repository

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
)

type SyncQueueRepository struct {
	DBConn *gorm.DB
}

func (r *SyncQueueRepository) Create(q *entity.SyncQueue) error {
	return r.DBConn.Create(q).Error
}

func (r *SyncQueueRepository) Update(q *entity.SyncQueue) error {
	return r.DBConn.Save(q).Error
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
	var existing entity.SyncQueue

	err := r.DBConn.
		Where("spreadsheet_id = ? AND sheet_name = ? AND form_notes = ?", q.SpreadsheetID, q.SheetName, q.FormNotes).
		First(&existing).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Chưa có thì tạo mới
			return r.DBConn.Create(q).Error
		}
		// Có lỗi truy vấn
		return err
	}

	// Đã có → cập nhật các field liên quan
	existing.LastSubmissionID = q.LastSubmissionID
	existing.LastSubmittedAt = q.LastSubmittedAt
	existing.Status = q.Status
	existing.UpdatedAt = q.UpdatedAt
	existing.FormNotes = q.FormNotes

	return r.DBConn.Save(&existing).Error
}

func (r *SyncQueueRepository) GetBySheetUrlAndSheetNameAndFormNotes(sheetUrl, sheetName string, formNotesJSON []byte) (*entity.SyncQueue, error) {
	var queue entity.SyncQueue

	err := r.DBConn.
		Where("sheet_url = ? AND sheet_name = ? AND form_notes = ?", sheetUrl, sheetName, formNotesJSON).
		First(&queue).Error

	if err != nil {
		return nil, err
	}
	return &queue, nil
}

func (r *SyncQueueRepository) UpdateStatusByID(id uint64, status value.SyncQueueStatus) error {
	return r.DBConn.Model(&entity.SyncQueue{}).
		Where("id = ?", id).
		Update("status", status).Error
}
