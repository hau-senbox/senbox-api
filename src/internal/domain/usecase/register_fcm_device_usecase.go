package usecase

import (
	firebase "firebase.google.com/go/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type RegisterFcmDeviceUseCase struct {
	FirebaseApp *firebase.App
	DBConn      *gorm.DB
	Repository  *repository.MobileDeviceRepository
}

func NewRegisterFcmDeviceUseCase(db *gorm.DB, app *firebase.App) *RegisterFcmDeviceUseCase {
	return &RegisterFcmDeviceUseCase{
		FirebaseApp: app,
		DBConn:      db,
		Repository:  &repository.MobileDeviceRepository{},
	}
}

func (r *RegisterFcmDeviceUseCase) Execute(req request.RegisterFCMRequest) error {
	_, err := r.Repository.Save(entity.SMobileDevice{
		DeviceId: req.DeviceId,
		FCMToken: req.DeviceToken,
		Type:     req.Type,
	}, r.DBConn)

	if err != nil {
		log.Error("RegisterFcmDeviceUseCase.Execute failed to save FCM device ", err)
	}

	return err
}
