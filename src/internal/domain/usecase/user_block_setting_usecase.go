package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"time"
)

type UserBlockSettingUsecase struct {
	Repo *repository.UserBlockSettingRepository
}

func NewUserBlockSettingUsecase(repo *repository.UserBlockSettingRepository) *UserBlockSettingUsecase {
	return &UserBlockSettingUsecase{Repo: repo}
}

// GetByUserID returns user block setting by user ID
func (uc *UserBlockSettingUsecase) GetByUserID(userID string) (*response.UserBlockSettingResponse, error) {
	setting, err := uc.Repo.GetByUserID(userID)
	if err != nil {
		return nil, err
	}

	if setting == nil {
		return nil, nil
	}

	return &response.UserBlockSettingResponse{
		ID:              setting.ID,
		UserID:          setting.UserID,
		IsDeactive:      setting.IsDeactive,
		IsViewMessage:   setting.IsViewMessage,
		MessageBox:      setting.MessageBox,
		MessageDeactive: setting.MessageDeactive,
		CreatedAt:       setting.CreatedAt,
		UpdatedAt:       setting.UpdatedAt,
	}, nil
}

// Create a new block setting
func (uc *UserBlockSettingUsecase) Create(req request.UserBlockSettingRequest) error {
	setting := &entity.UserBlockSetting{
		UserID:          req.UserID,
		IsDeactive:      req.IsDeactive,
		IsViewMessage:   req.IsViewMessage,
		MessageBox:      req.MessageBox,
		MessageDeactive: req.MessageDeactive,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	return uc.Repo.Create(setting)
}

// Update existing block setting
func (uc *UserBlockSettingUsecase) Update(id int, req request.UserBlockSettingRequest) error {
	setting, err := uc.Repo.GetByID(id)
	if err != nil {
		return err
	}

	setting.IsDeactive = req.IsDeactive
	setting.IsViewMessage = req.IsViewMessage
	setting.MessageBox = req.MessageBox
	setting.MessageDeactive = req.MessageDeactive
	setting.UpdatedAt = time.Now()

	return uc.Repo.Update(setting)
}

// Upsert (create if not exists, otherwise update)
func (uc *UserBlockSettingUsecase) Upsert(req request.UserBlockSettingRequest) error {
	setting, err := uc.Repo.GetByUserID(req.UserID)
	if err != nil {
		return err
	}

	if setting == nil {
		return uc.Create(req)
	}

	setting.IsDeactive = req.IsDeactive
	setting.IsViewMessage = req.IsViewMessage
	setting.MessageBox = req.MessageBox
	setting.MessageDeactive = req.MessageDeactive
	setting.UpdatedAt = time.Now()

	return uc.Repo.Update(setting)
}

// Delete block setting by ID
func (uc *UserBlockSettingUsecase) Delete(id int) error {
	return uc.Repo.Delete(id)
}
