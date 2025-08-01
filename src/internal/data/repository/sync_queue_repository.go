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
