package usecase

import (
	"errors"
	"sen-global-api/helper"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	"gorm.io/gorm"
)

type AppConfigUseCase struct {
	Repo *repository.AppConfigRepository
}

// Lấy tất cả configs
func (uc *AppConfigUseCase) GetAll() ([]entity.AppConfig, error) {
	return uc.Repo.GetAll()
}

// Upload: nếu tồn tại key thì update, nếu không thì create
func (uc *AppConfigUseCase) Upload(req request.UploadAppConfigRequest) error {
	// convert request → entity
	config := entity.AppConfig{
		Key:   req.Key,
		Value: helper.AppConfigValueToJSON(req.Value),
	}

	existing, err := uc.Repo.GetByKey(req.Key)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	if existing != nil && existing.ID != 0 {
		// update value
		existing.Value = config.Value
		return uc.Repo.Update(existing)
	}

	// create mới
	return uc.Repo.Create(&config)
}
