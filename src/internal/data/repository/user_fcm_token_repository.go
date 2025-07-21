package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"time"

	"gorm.io/gorm"
)

type UserTokenFCMRepository struct {
	DBConn *gorm.DB
}

func NewUserTokenFCMRepository(dbConn *gorm.DB) *UserTokenFCMRepository {
	return &UserTokenFCMRepository{DBConn: dbConn}
}

func (receiver *UserTokenFCMRepository) FindByDeviceID(userID, deviceID string) (*entity.SUserFCMToken, error) {

	var userFCMToken entity.SUserFCMToken

	err := receiver.DBConn.Model(&entity.SUserFCMToken{}).Where("user_id = ? AND device_id = ?", userID, deviceID).First(&userFCMToken).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &userFCMToken, nil

}

func (receiver *UserTokenFCMRepository) CreateToken(data *entity.SUserFCMToken) error {
	return receiver.DBConn.Model(&entity.SUserFCMToken{}).Create(data).Error
}

func (receiver *UserTokenFCMRepository) UpdateToken(data *entity.SUserFCMToken) error {
	return receiver.DBConn.Model(&entity.SUserFCMToken{}).
		Where("user_id = ? AND device_id = ?", data.UserID, data.DeviceID).
		UpdateColumns(map[string]interface{}{
			"fcm_token":  data.FCMToken,
			"updated_at": time.Now(),
			"is_active":  true,
		}).Error
}

func (receiver *UserTokenFCMRepository) FindByUserID(userID string) ([]*entity.SUserFCMToken, error) {
	
	var tokens []*entity.SUserFCMToken
	
	err := receiver.DBConn.Model(&entity.SUserFCMToken{}).Where("user_id = ?", userID).Find(&tokens).Error
	if err != nil {
		return nil, err
	}

	return tokens, nil
}