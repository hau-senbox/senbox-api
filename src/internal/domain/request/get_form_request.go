package request

type GetFormRequest struct {
	QrCode   string `json:"qr_code" binding:"required"`
	DeviceID string `json:"device_id" binding:"required"`
}
