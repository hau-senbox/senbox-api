package usecase

import (
	"errors"
	"sen-global-api/config"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type FindDeviceFromRequestCase struct {
	*repository.DeviceRepository
	*repository.SessionRepository
}

func (receiver *FindDeviceFromRequestCase) FindDevice(context *gin.Context) (*entity.SDevice, error) {
	authorization := context.GetHeader("Authorization")
	if authorization == "" {
		return nil, errors.New("no authorization header")
	}

	if len(authorization) == 0 {
		return nil, errors.New("no authorization header")
	}

	tokenString := strings.Split(authorization, " ")[1]

	deviceId, err := receiver.ExtractDeviceIdFromToken(tokenString)
	if err != nil || deviceId == nil {
		return nil, errors.New("invalid token")
	}

	device, err := receiver.FindDeviceById(*deviceId)
	if err != nil {
		return nil, err
	}

	return device, nil
}

func NewFindDeviceFromRequestCase(cfg config.AppConfig, dbConn *gorm.DB) *FindDeviceFromRequestCase {
	return &FindDeviceFromRequestCase{
		DeviceRepository: &repository.DeviceRepository{
			DBConn:                      dbConn,
			DefaultRequestPageSize:      cfg.DefaultRequestPageSize,
			DefaultOutputSpreadsheetUrl: cfg.OutputSpreadsheetUrl,
		},
		SessionRepository: &repository.SessionRepository{
			AuthorizeEncryptKey:   cfg.AuthorizeEncryptKey,
			TokenExpireTimeInHour: time.Duration(cfg.TokenExpireDurationInHour),
		},
	}
}
