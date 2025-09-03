package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
)

type UserSettingRepository struct {
	DBConn *gorm.DB
}

// Create user setting
func (r *UserSettingRepository) Create(setting *entity.UserSetting) error {
	return r.DBConn.Create(setting).Error
}

// Get by ID
func (r *UserSettingRepository) GetByID(id uint) (*entity.UserSetting, error) {
	var setting entity.UserSetting
	if err := r.DBConn.First(&setting, id).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// Get by OwnerID + Key
func (r *UserSettingRepository) GetByOwnerAndKey(ownerID, key string) (*entity.UserSetting, error) {
	var setting entity.UserSetting
	if err := r.DBConn.Where("owner_id = ? AND `key` = ?", ownerID, key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// Update value
func (r *UserSettingRepository) Update(setting *entity.UserSetting) error {
	return r.DBConn.Save(setting).Error
}

// Delete by ID
func (r *UserSettingRepository) Delete(id uint) error {
	result := r.DBConn.Delete(&entity.UserSetting{}, id)
	if result.RowsAffected == 0 {
		return errors.New("no record found to delete")
	}
	return result.Error
}

func (r *UserSettingRepository) GetByOwner(ownerID string) ([]*entity.UserSetting, error) {
	var settings []*entity.UserSetting
	if err := r.DBConn.Where("owner_id = ?", ownerID).Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

func (r *UserSettingRepository) GetLoginDeviceLimit(ownerID string) (*entity.UserSetting, error) {
	var setting entity.UserSetting
	if err := r.DBConn.Where("owner_id = ? AND `key` = ?", ownerID, string(value.UserSettingLoginDeviceLimit)).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}
