package request

import "sen-global-api/internal/domain/value"

type RegisterFCMRequest struct {
	DeviceID    string           `json:"device_id" binding:"required"`
	DeviceToken string           `json:"device_token" binding:"required"`
	Type        value.DeviceType `json:"type" binding:"required"`
}
