package request

type GetRedirectUrlByQRCodeRequest struct {
	QRCode string `form:"qr_code"`
}
