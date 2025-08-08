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
