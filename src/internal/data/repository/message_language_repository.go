package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type MessageLanguageRepository struct {
	DBConn *gorm.DB
}

func NewMessageLanguageRepository(dbConn *gorm.DB) *MessageLanguageRepository {
	return &MessageLanguageRepository{
		DBConn: dbConn,
	}
}

func (r *MessageLanguageRepository) GetAll() ([]entity.MessageLanguage, error) {
	var messages []entity.MessageLanguage
	if err := r.DBConn.Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageLanguageRepository) GetByTypeAndLanguage(typeStr string, languageID uint) ([]entity.MessageLanguage, error) {
	var messages []entity.MessageLanguage
	if err := r.DBConn.Where("type = ? AND language_id = ?", typeStr, languageID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageLanguageRepository) Create(message *entity.MessageLanguage) error {
	return r.DBConn.Create(message).Error
}

func (r *MessageLanguageRepository) Update(message *entity.MessageLanguage) error {
	return r.DBConn.Save(message).Error
}

func (r *MessageLanguageRepository) Delete(id int) error {
	return r.DBConn.Delete(&entity.MessageLanguage{}, "id = ?", id).Error
}

func (r *MessageLanguageRepository) GetByTypeAndKeyAndLanguage(typeStr string, key string, languageID uint) ([]entity.MessageLanguage, error) {
	var messages []entity.MessageLanguage
	if err := r.DBConn.Where("type = ? AND language_id = ? AND `key` = ?", typeStr, languageID, key).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageLanguageRepository) GetByTypeAndKeyAndLanguageAndTypeID(typeStr string, key string, languageID uint, typeID string) (*entity.MessageLanguage, error) {
	var messages *entity.MessageLanguage
	if err := r.DBConn.Where("type = ? AND language_id = ? AND `key` = ? AND type_id = ?", typeStr, languageID, key, typeID).First(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}

func (r *MessageLanguageRepository) GetByTypeAndTypeID(typeStr string, typeID string) ([]entity.MessageLanguage, error) {
	var messages []entity.MessageLanguage
	if err := r.DBConn.Where("type = ? AND type_id = ?", typeStr, typeID).Find(&messages).Error; err != nil {
		return nil, err
	}
	return messages, nil
}
