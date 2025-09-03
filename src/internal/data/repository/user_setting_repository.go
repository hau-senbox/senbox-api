package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type UserSettingRepository struct {
	DB *gorm.DB
}

// Create user setting
func (r *UserSettingRepository) Create(setting *entity.UserSetting) error {
	return r.DB.Create(setting).Error
}

// Get by ID
func (r *UserSettingRepository) GetByID(id uint) (*entity.UserSetting, error) {
	var setting entity.UserSetting
	if err := r.DB.First(&setting, id).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// Get by OwnerID + Key
func (r *UserSettingRepository) GetByOwnerAndKey(ownerID, key string) (*entity.UserSetting, error) {
	var setting entity.UserSetting
	if err := r.DB.Where("owner_id = ? AND `key` = ?", ownerID, key).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// Update value
func (r *UserSettingRepository) Update(setting *entity.UserSetting) error {
	return r.DB.Save(setting).Error
}

// Delete by ID
func (r *UserSettingRepository) Delete(id uint) error {
	result := r.DB.Delete(&entity.UserSetting{}, id)
	if result.RowsAffected == 0 {
		return errors.New("no record found to delete")
	}
	return result.Error
}
