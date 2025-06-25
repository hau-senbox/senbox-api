package request

type SendNotificationRequest struct {
	DeviceID string `json:"device_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Message  string
}
