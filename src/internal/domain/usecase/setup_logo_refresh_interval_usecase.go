package usecase

import (
	"context"
	"firebase.google.com/go/v4/messaging"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"sen-global-api/pkg/monitor"
	"strconv"
)

func SetupLogoRefreshIntervalUseCase(req request.SetupLogoRefreshIntervalRequest) error {
	repo := repository.NewSettingRepository(DBConn)

	if req.Interval > 0 {
		err := repo.UpdateLogoRefreshInterval(req.Interval)
		if err != nil {
			return err
		}
		go func() {
			announceLogoFreshUpdatedInterval(req.Interval)
		}()
	}

	if req.Title != "" {
		err := repo.UpdateLogoRefreshTitle(req.Title)
		if err != nil {
			return err
		}
	}

	return nil
}

func announceLogoFreshUpdatedInterval(interval uint64) {
	monitor.SendMessageViaTelegram("Logo refresh interval updated to ", strconv.FormatUint(interval, 10))
	ctx := context.Background()
	msgApp, err := FirebaseApp.Messaging(ctx)
	if err != nil {
		monitor.SendMessageViaTelegram("[ERROR][INFORM LOGO REFRESH INTERVAL] Cannot initialize Messaging App ", err.Error())
		return
	}

	msg := &messaging.Message{
		Topic: string(value.FcmTopicsGeneral),
		Android: &messaging.AndroidConfig{
			Priority: "high",
			Notification: &messaging.AndroidNotification{
				Icon: "ic_launcher",
			},
			Data: map[string]string{
				"interval": strconv.FormatUint(interval, 10),
				"type":     string(value.NotificationType_LogoRefreshIntervalChanged),
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					ContentAvailable: true,
				},
				CustomData: map[string]interface{}{
					"interval": strconv.FormatUint(interval, 10),
					"type":     value.NotificationType_LogoRefreshIntervalChanged,
				},
			},
		},
	}

	_, err = msgApp.Send(ctx, msg)
	if err != nil {
		monitor.SendMessageViaTelegram("[ERROR][INFORM LOGO REFRESH INTERVAL] Cannot send notification", " error ", err.Error())
	}
}

func GetLogoRefreshIntervalUseCase() (entity.SSetting, error) {
	repo := repository.NewSettingRepository(DBConn)
	return repo.GetLogoRefreshInterval()
}
