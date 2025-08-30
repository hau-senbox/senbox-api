package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type UserDevicesLoginRepository struct {
	DBConn *gorm.DB
}

// Create lưu thông tin login mới
func (r *UserDevicesLoginRepository) Create(userDevice *entity.UserDevicesLogin) error {
	return r.DBConn.Create(userDevice).Error
}

// GetByUserAndDevice lấy bản ghi theo UserID và DeviceID
func (r *UserDevicesLoginRepository) GetByUserAndDevice(userID, deviceID string) (*entity.UserDevicesLogin, error) {
	var record entity.UserDevicesLogin
	err := r.DBConn.Where("user_id = ? AND device_id = ?", userID, deviceID).First(&record).Error
	if err != nil {
		return nil, err
	}
	return &record, nil
}

// GetByUser lấy danh sách device login theo UserID
func (r *UserDevicesLoginRepository) GetByUser(userID string) ([]entity.UserDevicesLogin, error) {
	var records []entity.UserDevicesLogin
	err := r.DBConn.Where("user_id = ?", userID).Find(&records).Error
	if err != nil {
		return nil, err
	}
	return records, nil
}

// DeleteByUserAndDevice xóa login theo UserID + DeviceID
func (r *UserDevicesLoginRepository) DeleteByUserAndDevice(userID string, deviceID string) error {
	return r.DBConn.Where("user_id = ? AND device_id = ?", userID, deviceID).Delete(&entity.UserDevicesLogin{}).Error
}

// DeleteByUser xóa tất cả login của user
func (r *UserDevicesLoginRepository) DeleteByUser(userID string) error {
	return r.DBConn.Where("user_id = ?", userID).Delete(&entity.UserDevicesLogin{}).Error
}

// CountDevicesByUserExcludeDevice đếm số lượng device login theo UserID
// nhưng loại trừ 1 DeviceID cụ thể
func (r *UserDevicesLoginRepository) CountDevicesByUserExcludeDevice(userID, excludeDeviceID string) (int64, error) {
	var count int64
	err := r.DBConn.Model(&entity.UserDevicesLogin{}).
		Where("user_id = ? AND device_id <> ?", userID, excludeDeviceID).
		Count(&count).Error
	if err != nil {
		return 0, err
	}
	return count, nil
}

// Update cập nhật bản ghi UserDevicesLogin
func (r *UserDevicesLoginRepository) Update(userDevice *entity.UserDevicesLogin) error {
	return r.DBConn.Save(userDevice).Error
}
