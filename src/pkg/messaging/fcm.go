package messaging

import (
	"context"
	"fmt"
	"os"
	"sen-global-api/config"
	"sen-global-api/internal/domain/value"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/messaging"
	log "github.com/sirupsen/logrus"
	"google.golang.org/api/option"
)

func NewFirebaseApp(cfg config.AppConfig) (*firebase.App, error) {
	credentialsInByte, err := os.ReadFile(cfg.Messaging.ServiceAccount)
	if err != nil {
		return nil, err
	}

	opt := option.WithCredentialsJSON(credentialsInByte)
	otherApp, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, err
	}

	return otherApp, err
}

type NotificationParams struct {
	Title       string
	Message     string
	DeviceToken string
	Type        value.NotificationType
}

func SendNotification(app *firebase.App, params NotificationParams) error {
	ctx := context.Background()
	msgApp, err := app.Messaging(ctx)
	if err != nil {
		return fmt.Errorf("cannot initialize Messaging App %s", err.Error())
	}

	msg := &messaging.Message{
		Notification: &messaging.Notification{
			Title: params.Title,
			Body:  params.Message,
		},
		Token: params.DeviceToken,
		Android: &messaging.AndroidConfig{
			Priority:    "high",
			CollapseKey: string(params.Type),
			Notification: &messaging.AndroidNotification{
				Icon:  "ic_launcher",
				Title: params.Title,
				Body:  params.Message,
			},
			Data: map[string]string{
				"type": string(params.Type),
			},
		},
		APNS: &messaging.APNSConfig{
			Headers: map[string]string{
				"apns-priority": "10",
			},
			Payload: &messaging.APNSPayload{
				Aps: &messaging.Aps{
					Alert: &messaging.ApsAlert{
						Title: params.Title,
						Body:  params.Message,
					},
				},
				CustomData: map[string]interface{}{
					"type": string(params.Type),
				},
			},
		},
	}

	_, err = msgApp.Send(ctx, msg)
	if err != nil {
		log.Errorf("FCM failed to send message to device token %s error %s", params.DeviceToken, err.Error())
	}

	return err
}
