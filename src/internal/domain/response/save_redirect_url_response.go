package response

import "time"

type SaveRedirectUrlResponseData struct {
	Id        uint64    `json:"id" binding:"required"`
	QRCode    string    `json:"qr_code" binding:"required"`
	TargetUrl string    `json:"target_url" binding:"required"`
	Password  *string   `json:"password"`
	CreatedAt time.Time `json:"created_at" binding:"required"`
	UpdatedAt time.Time `json:"updated_at" binding:"required"`
}

type SaveRedirectUrlResponse struct {
	Data SaveRedirectUrlResponseData `json:"data" binding:"required"`
}
