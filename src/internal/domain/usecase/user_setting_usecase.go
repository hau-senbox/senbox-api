package usecase

import (
	"encoding/json"
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
)

type UserSettingUseCase struct {
	Repo *repository.UserSettingRepository
}

// Get setting by ID
func (uc *UserSettingUseCase) GetUserSettingByID(id uint) (*entity.UserSetting, error) {
	return uc.Repo.GetByID(id)
}

// Get setting by OwnerID + Key
func (uc *UserSettingUseCase) GetUserSetting(ownerID, key string) (*entity.UserSetting, error) {
	return uc.Repo.GetByOwnerAndKey(ownerID, key)
}

// Update setting value
func (uc *UserSettingUseCase) UpdateUserSetting(ownerID, key string, value []byte) (*entity.UserSetting, error) {
	setting, err := uc.Repo.GetByOwnerAndKey(ownerID, key)
	if err != nil {
		return nil, errors.New("setting not found")
	}

	setting.Value = value
	if err := uc.Repo.Update(setting); err != nil {
		return nil, err
	}
	return setting, nil
}

// Delete setting by ID
func (uc *UserSettingUseCase) DeleteUserSetting(id uint) error {
	return uc.Repo.Delete(id)
}

func (uc *UserSettingUseCase) UploadUserSetting(req request.UploadUserSettingRequest) (*entity.UserSetting, error) {
	// Validate OwnerRole
	if !value.OwnerRole(req.OwnerRole).IsValid() {
		return nil, errors.New("invalid owner role")
	}

	// Validate Key
	if !value.UserSettingKey(req.Key).IsValid() {
		return nil, errors.New("invalid user setting key")
	}

	valueBytes, err := json.Marshal(req.Value)
	if err != nil {
		return nil, errors.New("invalid value format")
	}

	// Check setting exists
	setting, err := uc.Repo.GetByOwnerAndKey(req.OwnerID, req.Key)
	if err != nil {
		// nếu lỗi là record not found -> tạo mới
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newSetting := &entity.UserSetting{
				OwnerID:   req.OwnerID,
				OwnerRole: value.OwnerRole(req.OwnerRole),
				Key:       value.UserSettingKey(req.Key),
				Value:     valueBytes,
			}
			if err := uc.Repo.Create(newSetting); err != nil {
				return nil, err
			}
			return newSetting, nil
		}
		// nếu lỗi khác -> return
		return nil, err
	}

	// Update nếu đã tồn tại
	setting.Value = valueBytes
	if err := uc.Repo.Update(setting); err != nil {
		return nil, err
	}

	return setting, nil
}

func (uc *UserSettingUseCase) GetByOwner(ownerID string) (*response.UserSettingResponse, error) {
	settings, err := uc.Repo.GetByOwner(ownerID)
	if err != nil {
		return nil, err
	}

	return mapper.ToUserSettingResponse(settings), nil
}
