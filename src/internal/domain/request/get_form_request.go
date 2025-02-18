package request

type GetFormRequest struct {
	QrCode   string `json:"qr_code" binding:"required"`
	DeviceId string `json:"device_id" binding:"required"`
}
