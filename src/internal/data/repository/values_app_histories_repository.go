package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type ValuesAppHistoriesRepository struct {
	DBConn *gorm.DB
}

func (r *ValuesAppHistoriesRepository) Create(histories *entity.ValuesAppHistories) error {
	return r.DBConn.Create(histories).Error
}

func (r *ValuesAppHistoriesRepository) Update(histories *entity.ValuesAppHistories) error {
	return r.DBConn.Save(histories).Error
}

func (r *ValuesAppHistoriesRepository) GetByDeviceID(deviceID string) ([]*entity.ValuesAppHistories, error) {
	var histories []*entity.ValuesAppHistories
	if err := r.DBConn.Where("device_id = ?", deviceID).Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}

func (r *ValuesAppHistoriesRepository) GetByDeviceIDAndOrgID(deviceID string, orgID string) ([]*entity.ValuesAppHistories, error) {
	var histories []*entity.ValuesAppHistories
	if err := r.DBConn.Where("device_id = ? AND value2 = ?", deviceID, orgID).Find(&histories).Error; err != nil {
		return nil, err
	}
	return histories, nil
}
