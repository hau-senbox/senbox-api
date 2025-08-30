package usecase

import (
	"errors"
	"time"

	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
)

type ManageUserLoginUseCase struct {
	UserDevicesLoginRepository *repository.UserDevicesLoginRepository
}

// ManageUserDeviceLogin tạo mới login cho user
// điều kiện: 1 user chỉ được login tối đa 2 device
func (uc *ManageUserLoginUseCase) ManageUserDeviceLogin(userID, deviceID string) error {
	// Kiểm tra device đã tồn tại cho user chưa
	existing, err := uc.UserDevicesLoginRepository.GetByUserAndDevice(userID, deviceID)
	if err == nil && existing != nil {
		// Đã tồn tại -> update LoginAt để refresh thời gian login
		existing.LoginAt = time.Now()
		return uc.UserDevicesLoginRepository.Update(existing)
	}

	// Đếm số device hiện tại (ngoại trừ deviceID này)
	count, err := uc.UserDevicesLoginRepository.CountDevicesByUserExcludeDevice(userID, deviceID)
	if err != nil {
		return err
	}
	if count >= 2 {
		return errors.New("user has already logged in with the maximum of 2 devices")
	}

	// add device ?

	// Thêm mới
	newLogin := &entity.UserDevicesLogin{
		UserID:   userID,
		DeviceID: deviceID,
		LoginAt:  time.Now(),
	}
	return uc.UserDevicesLoginRepository.Create(newLogin)
}

func (uc *ManageUserLoginUseCase) ManageUserDeviceLogout(userID string, deviceID string) error {
	// xoa device user login
	err := uc.UserDevicesLoginRepository.DeleteByUserAndDevice(userID, deviceID)
	if err != nil {
		return err
	}
	return nil
}
