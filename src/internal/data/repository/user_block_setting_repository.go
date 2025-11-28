package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type UserBlockSettingRepository struct {
	DBConn *gorm.DB
}

func (r *UserBlockSettingRepository) Create(setting *entity.UserBlockSetting) error {
	return r.DBConn.Create(setting).Error
}

func (r *UserBlockSettingRepository) Update(setting *entity.UserBlockSetting) error {
	return r.DBConn.Save(setting).Error
}

func (r *UserBlockSettingRepository) GetByUserID(userID string) (*entity.UserBlockSetting, error) {
	var setting entity.UserBlockSetting
	err := r.DBConn.Where("user_id = ?", userID).First(&setting).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *UserBlockSettingRepository) GetByID(id int) (*entity.UserBlockSetting, error) {
	var setting entity.UserBlockSetting
	err := r.DBConn.First(&setting, id).Error
	if err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *UserBlockSettingRepository) Delete(id int) error {
	return r.DBConn.Delete(&entity.UserBlockSetting{}, id).Error
}

func (r *UserBlockSettingRepository) GetIsDeactiveByUserID(userID string) (bool, error) {
	var isDeactive bool
	err := r.DBConn.Model(&entity.UserBlockSetting{}).
		Select("is_deactive").
		Where("user_id = ?", userID).
		Scan(&isDeactive).Error
	if err != nil {
		return false, err
	}
	return isDeactive, nil
}

func (r *UserBlockSettingRepository) OnIsNeedToUpdate() error {
	return r.DBConn.Model(&entity.UserBlockSetting{}).
		Where("is_need_to_update = ?", false).
		Update("is_need_to_update", true).Error
}

func (r *UserBlockSettingRepository) OffIsNeedToUpdate(userID string) error {
	return r.DBConn.Model(&entity.UserBlockSetting{}).
		Where("user_id = ?", userID).
		Update("is_need_to_update", false).Error
}

func (r *UserBlockSettingRepository) GetAll() ([]entity.UserBlockSetting, error) {
	var userBlockSettings []entity.UserBlockSetting
	err := r.DBConn.Find(&userBlockSettings).Error
	if err != nil {
		return nil, err
	}
	return userBlockSettings, nil
}
