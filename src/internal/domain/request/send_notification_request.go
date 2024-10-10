package request

type SendNotificationRequest struct {
	DeviceId string `json:"device_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Message  string
}
