package request

import "time"

type SubmitFormRequest struct {
	QRCode   string    `json:"qr_code" binding:"required"`
	Answers  []Answer  `json:"answers" binding:"required"`
	OpenedAt time.Time `json:"opened_at,default=now()"`
	DeviceId string    `json:"device_id"`
}
