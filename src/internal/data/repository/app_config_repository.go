package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type AppConfigRepository struct {
	DBConn *gorm.DB
}

func (r *AppConfigRepository) GetByKey(key string) (*entity.AppConfig, error) {
	var config entity.AppConfig
	if err := r.DBConn.Where("`key` = ?", key).First(&config).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *AppConfigRepository) GetAll() ([]entity.AppConfig, error) {
	var configs []entity.AppConfig
	if err := r.DBConn.Find(&configs).Error; err != nil {
		return nil, err
	}
	return configs, nil
}

func (r *AppConfigRepository) Create(config *entity.AppConfig) error {
	return r.DBConn.Create(config).Error
}

func (r *AppConfigRepository) Update(config *entity.AppConfig) error {
	return r.DBConn.Save(config).Error
}

func (r *AppConfigRepository) Delete(id uint) error {
	return r.DBConn.Delete(&entity.AppConfig{}, id).Error
}
