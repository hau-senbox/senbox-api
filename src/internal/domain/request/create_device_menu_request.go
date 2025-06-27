package request

import "sen-global-api/internal/domain/entity/menu"

type CreateDeviceMenuRequest struct {
	DeviceID   string                       `json:"device_id" binding:"required"`
	Direction  menu.Direction               `json:"direction" binding:"required"`
	Components []CreateMenuComponentRequest `json:"components" binding:"required"`
}
