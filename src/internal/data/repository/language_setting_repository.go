package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type LanguageSettingRepository struct {
	DBConn *gorm.DB
}

func (r *LanguageSettingRepository) Create(setting *entity.LanguageSetting) error {
	return r.DBConn.Create(setting).Error
}

func (r *LanguageSettingRepository) Update(setting *entity.LanguageSetting) error {
	return r.DBConn.Save(setting).Error
}

func (r *LanguageSettingRepository) GetByID(id uint) (*entity.LanguageSetting, error) {
	var setting entity.LanguageSetting
	if err := r.DBConn.First(&setting, id).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

func (r *LanguageSettingRepository) GetAll() ([]entity.LanguageSetting, error) {
	var settings []entity.LanguageSetting
	err := r.DBConn.Find(&settings).Error
	return settings, err
}

func (r *LanguageSettingRepository) DeleteAll() error {
	return r.DBConn.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.LanguageSetting{}).Error
}

func (r *LanguageSettingRepository) Delete(id uint) error {
	return r.DBConn.Delete(&entity.LanguageSetting{}, id).Error
}

func (r *LanguageSettingRepository) DeleteMany(ids []uint) error {
	return r.DBConn.Where("id IN ?", ids).Delete(&entity.LanguageSetting{}).Error
}
