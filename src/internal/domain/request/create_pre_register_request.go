package request

type CreatePreRegisterRequest struct {
	Email      string `json:"email" binding:"required"`
	DeviceID   string `json:"device_id" binding:"required"`
	DeviceName string `json:"device_name"`
	FormQr     string `json:"form_qr" binding:"required"`
}
