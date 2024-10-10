package request

type SaveRedirectUrlRequest struct {
	QRCode    string `json:"qr_code" binding:"required"`
	TargetUrl string `json:"target_url" binding:"required"`
	Password  string `json:"password"`
}
