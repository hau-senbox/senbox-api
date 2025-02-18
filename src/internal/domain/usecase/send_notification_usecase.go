package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/messaging"
	"sen-global-api/pkg/monitor"

	firebase "firebase.google.com/go/v4"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var FirebaseApp *firebase.App

type SendNotificationUseCase struct {
	FirebaseApp *firebase.App
	repository  *repository.MobileDeviceRepository
	DB          *gorm.DB
}

func NewSendNotificationUseCase(conn *gorm.DB, app *firebase.App) *SendNotificationUseCase {
	return &SendNotificationUseCase{
		FirebaseApp: app,
		repository:  repository.NewMobileDeviceRepository(),
		DB:          conn,
	}
}

type SendNotificationParams struct {
	DeviceToken string
}

func (receiver *SendNotificationUseCase) Execute(params request.SendNotificationRequest) error {
	md, err := receiver.repository.FindByDeviceID(params.DeviceId, receiver.DB)

	if err != nil {
		return err
	}

	noti := messaging.NotificationParams{
		Title:       params.Title,
		Message:     "",
		DeviceToken: md.FCMToken,
		Type:        value.NotificationType_NewFormSubmit,
	}

	err = messaging.SendNotification(receiver.FirebaseApp, noti)
	if err != nil {
		log.Error("Failed to send notification ", err)
		monitor.SendMessageViaTelegram("Failed to send notification for a top buttons: ", err.Error())
	}

	return err
}
