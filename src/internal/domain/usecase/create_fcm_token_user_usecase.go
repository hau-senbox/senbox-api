package usecase

import (
	"errors"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"time"

)

type CreateUserTokenFCMUseCase struct {
	UserTokenFCMRepository *repository.UserTokenFCMRepository
}

func (receiver *CreateUserTokenFCMUseCase) CreateToken(userID, deviceID, token string) error {

	var data *entity.SUserFCMToken

	if userID == "" {
		return errors.New("user id is required")
	}

	if deviceID == "" {
		return errors.New("device id is required")
	}

	if token == "" {
		return errors.New("token is required")
	}

	data, err := receiver.UserTokenFCMRepository.FindByDeviceID(userID, deviceID)
	if err != nil {
		return err
	}

	if data != nil {
		data.FCMToken = token
		return receiver.UserTokenFCMRepository.UpdateToken(data)
	} else {
		data = &entity.SUserFCMToken{
			UserID:    userID,
			DeviceID:  deviceID,
			FCMToken:  token,
			IsActive:  true,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		return receiver.UserTokenFCMRepository.CreateToken(data)
	}

}
