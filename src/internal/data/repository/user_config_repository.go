package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type UserConfigRepository struct {
	DBConn *gorm.DB
}

func NewUserConfigRepository(dbConn *gorm.DB) *UserConfigRepository {
	return &UserConfigRepository{DBConn: dbConn}
}

func (receiver *UserConfigRepository) GetByID(id uint) (*entity.SUserConfig, error) {
	var userConfig entity.SUserConfig
	err := receiver.DBConn.Where("id = ?", id).First(&userConfig).Error
	if err != nil {
		log.Error("UserConfigRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get user config")
	}
	return &userConfig, nil
}

func (receiver *UserConfigRepository) CreateUserConfig(req request.CreateUserConfigRequest) (*int64, error) {
	var userConfig entity.SUserConfig = entity.SUserConfig{
		TopButtonConfig:      req.TopButtonConfig,
		StudentOutputSheetId: req.StudentOutputSheetId,
		TeacherOutputSheetId: req.TeacherOutputSheetId,
	}

	result := receiver.DBConn.Create(&userConfig)

	if result.Error != nil {
		log.Error("UserConfigRepository.CreateUserConfig: " + result.Error.Error())
		return nil, errors.New("failed to create user config")
	}

	return &userConfig.ID, nil
}

func (receiver *UserConfigRepository) UpdateUserConfig(req request.UpdateUserConfigRequest) error {
	updateResult := receiver.DBConn.Model(&entity.SUserConfig{}).Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"top_button_config":       req.TopButtonConfig,
			"profile_picture_url":     req.ProfilePictureUrl,
			"student_output_sheet_id": req.StudentOutputSheetId,
			"teacher_output_sheet_id": req.TeacherOutputSheetId,
		})

	if updateResult.Error != nil {
		log.Error("UserConfigRepository.UpdateUserConfig: " + updateResult.Error.Error())
		return errors.New("failed to update user config")
	}

	return nil
}
