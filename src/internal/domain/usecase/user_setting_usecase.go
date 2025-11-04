package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/datatypes"
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

func (uc *UserSettingUseCase) UploadUserIsFirstLogin(ctx *gin.Context, request request.UploadUserIsFirstLoginRequest) error {

	// Check setting exists
	setting, err := uc.Repo.GetByOwnerAndKey(request.UserID, string(value.UserSettingIsFirstLogin))
	if err != nil {
		// nếu lỗi là record not found -> tạo mới
		if errors.Is(err, gorm.ErrRecordNotFound) {
			newSetting := &entity.UserSetting{
				OwnerID:   request.UserID,
				OwnerRole: value.OwnerRoleUser,
				Key:       value.UserSettingKey(value.UserSettingIsFirstLogin),
				Value:     datatypes.JSON([]byte(strconv.FormatBool(request.IsFirstLogin))),
			}
			if err := uc.Repo.Create(newSetting); err != nil {
				return err
			}
			return nil
		}
		// nếu lỗi khác -> return
		return err
	}

	// Update nếu đã tồn tại
	setting.Value = datatypes.JSON([]byte(strconv.FormatBool(request.IsFirstLogin)))
	if err := uc.Repo.Update(setting); err != nil {
		return err
	}
	return nil
}

func (uc *UserSettingUseCase) GetUserIsFirstLogin(userID string) (bool, error) {
	setting, err := uc.Repo.GetByOwnerAndKey(userID, string(value.UserSettingIsFirstLogin))
	if err != nil {
		return false, err
	}
	var isFirstLogin bool
	err = json.Unmarshal(setting.Value, &isFirstLogin)
	if err != nil {
		return false, errors.New("failed to unmarshal is first login")
	}
	return isFirstLogin, nil
}

func (uc *UserSettingUseCase) UploadUserWelcomeReminder(ctx *gin.Context, request request.UploadUserWelcomeReminderRequest) error {
	// Cộng thêm 3 ngày từ thời điểm hiện tại
	timeReminder := time.Now().AddDate(0, 0, 3).Format("2006-01-02 15:04:05")

	welcomeReminder := entity.UserSettingWelcomeReminder{
		IsEnabled:    request.IsEnabled,
		TimeReminder: timeReminder,
	}

	valueBytes, err := json.Marshal(welcomeReminder)
	if err != nil {
		return fmt.Errorf("failed to marshal reminder: %w", err)
	}

	// Kiểm tra xem setting đã tồn tại chưa
	setting, err := uc.Repo.GetByOwnerAndKey(request.UserID, string(value.UserSettingWelcomeReminder))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Tạo mới nếu chưa có
			newSetting := &entity.UserSetting{
				OwnerID:   request.UserID,
				OwnerRole: value.OwnerRoleUser,
				Key:       value.UserSettingKey(value.UserSettingWelcomeReminder),
				Value:     valueBytes,
			}
			return uc.Repo.Create(newSetting)
		}
		// Nếu lỗi khác thì trả về luôn
		return err
	}

	// Cập nhật nếu đã tồn tại
	setting.Value = valueBytes
	return uc.Repo.Update(setting)
}

func (uc *UserSettingUseCase) GetUserWelcomeReminder(userID string) (*entity.UserSettingWelcomeReminder, error) {
	setting, err := uc.Repo.GetByOwnerAndKey(userID, string(value.UserSettingWelcomeReminder))
	if err != nil {
		return nil, err
	}
	var welcomeReminder entity.UserSettingWelcomeReminder
	err = json.Unmarshal(setting.Value, &welcomeReminder)
	if err != nil {
		return nil, errors.New("failed to unmarshal welcome reminder")
	}
	return &welcomeReminder, nil
}
