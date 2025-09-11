package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type ValuesAppCurrentRepository struct {
	DBConn *gorm.DB
}

func (r *ValuesAppCurrentRepository) Create(log *entity.ValuesAppCurrent) error {
	return r.DBConn.Create(log).Error
}

func (r *ValuesAppCurrentRepository) Update(log *entity.ValuesAppCurrent) error {
	return r.DBConn.Save(log).Error
}

func (r *ValuesAppCurrentRepository) GetByID(id uint64) (*entity.ValuesAppCurrent, error) {
	var log entity.ValuesAppCurrent
	if err := r.DBConn.First(&log, id).Error; err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *ValuesAppCurrentRepository) GetAll() ([]*entity.ValuesAppCurrent, error) {
	var logs []*entity.ValuesAppCurrent
	if err := r.DBConn.Find(&logs).Error; err != nil {
		return nil, err
	}
	return logs, nil
}

func (r *ValuesAppCurrentRepository) FindByDeviceID(deviceID string) (*entity.ValuesAppCurrent, error) {
	var log entity.ValuesAppCurrent
	if err := r.DBConn.Where("device_id = ?", deviceID).First(&log).Error; err != nil {
		return nil, err
	}
	return &log, nil
}
