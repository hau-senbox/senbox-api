package repository

import (
	"sen-global-api/internal/domain/entity"

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
