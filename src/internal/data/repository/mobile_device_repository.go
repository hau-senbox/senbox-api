package repository

import (
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"sen-global-api/internal/domain/entity"
)

type MobileDeviceRepository struct {
}

func NewMobileDeviceRepository() *MobileDeviceRepository {
	return &MobileDeviceRepository{}
}

func (receiver *MobileDeviceRepository) Save(device entity.SMobileDevice, db *gorm.DB) (entity.SMobileDevice, error) {
	err := db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "device_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"fcm_token", "updated_at"}),
	}).Create(&device).Error

	return device, err
}

func (receiver *MobileDeviceRepository) FindByDeviceID(deviceId string, db *gorm.DB) (entity.SMobileDevice, error) {
	var d entity.SMobileDevice
	err := db.Where("device_id = ?", deviceId).First(&d).Error

	return d, err
}
