package request

type GetRedirectUrlByQRCodeRequest struct {
	QRCode string `form:"qrcode"`
}
