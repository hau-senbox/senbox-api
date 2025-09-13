package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type AccountsLogRepository struct {
	DBConn *gorm.DB
}

func (r *AccountsLogRepository) Create(accountsLog *entity.AccountsLog) error {
	return r.DBConn.Create(accountsLog).Error
}

func (r *AccountsLogRepository) Update(accountsLog *entity.AccountsLog) error {
	return r.DBConn.Save(accountsLog).Error
}

func (r *AccountsLogRepository) GetAll() ([]entity.AccountsLog, error) {
	var accountsLogs []entity.AccountsLog
	err := r.DBConn.Find(&accountsLogs).Error
	return accountsLogs, err
}

func (r *AccountsLogRepository) GetByUserID(userID string) ([]entity.AccountsLog, error) {
	var accountsLogs []entity.AccountsLog
	err := r.DBConn.Where("user_id = ?", userID).Find(&accountsLogs).Error
	return accountsLogs, err
}

func (r *AccountsLogRepository) GetByOrganizationID(organizationID string) ([]entity.AccountsLog, error) {
	var accountsLogs []entity.AccountsLog
	err := r.DBConn.Where("organization_id = ?", organizationID).Find(&accountsLogs).Error
	return accountsLogs, err
}

func (r *AccountsLogRepository) GetByDeviceID(deviceID string) ([]entity.AccountsLog, error) {
	var accountsLogs []entity.AccountsLog
	err := r.DBConn.Where("device_id = ?", deviceID).Find(&accountsLogs).Error
	return accountsLogs, err
}

func (r *AccountsLogRepository) GetByDeviceIDAndUserID(deviceID string, userID string) ([]entity.AccountsLog, error) {
	var accountsLogs []entity.AccountsLog
	err := r.DBConn.Where("device_id = ? AND user_id = ?", deviceID, userID).Find(&accountsLogs).Error
	return accountsLogs, err
}

func (r *AccountsLogRepository) GetByDeviceIDAndOrganizationID(deviceID string, organizationID string) ([]entity.AccountsLog, error) {
	var accountsLogs []entity.AccountsLog
	err := r.DBConn.Where("device_id = ? AND organization_id = ?", deviceID, organizationID).Find(&accountsLogs).Error
	return accountsLogs, err
}

func (r *AccountsLogRepository) GetByDeviceIDAndOrganizationIDAndUserID(deviceID string, organizationID string, userID string) ([]entity.AccountsLog, error) {
	var accountsLogs []entity.AccountsLog
	err := r.DBConn.Where("device_id = ? AND organization_id = ? AND user_id = ?", deviceID, organizationID, userID).Find(&accountsLogs).Error
	return accountsLogs, err
}
