package repository

import (
	"fmt"
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

func (r *LanguageSettingRepository) DeleteByIDs(tx *gorm.DB, ids []uint, compRepo *ComponentRepository) error {
	for _, id := range ids {
		// Lấy record language setting
		var langSetting entity.LanguageSetting
		if err := tx.First(&langSetting, id).Error; err != nil {
			return err
		}

		// Check tồn tại trong Component
		exist, err := compRepo.CheckExistLanguage(tx, langSetting.ID)
		if err != nil {
			return err
		}
		if exist {
			return fmt.Errorf("language setting ID %d is in use by components and cannot be deleted", langSetting.ID)
		}

		// Nếu không tồn tại trong component thì xóa
		if err := tx.Delete(&langSetting).Error; err != nil {
			return err
		}
	}
	return nil
}
