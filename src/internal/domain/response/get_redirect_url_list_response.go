package response

import "time"

type GetRedirectUrlListResponseData struct {
	Id           uint64    `json:"id" binding:"required"`
	QRCode       string    `json:"qr_code" binding:"required"`
	TargetUrl    string    `json:"target_url" binding:"required"`
	Password     *string   `json:"password" binding:"required"`
	Hint         string    `json:"hint" binding:"required"`
	HashPassword *string   `json:"hash_password" binding:"required"`
	CreatedAt    time.Time `json:"created_at" binding:"required"`
	UpdatedAt    time.Time `json:"updated_at" binding:"required"`
}

type GetRedirectUrlListResponse struct {
	Data   []GetRedirectUrlListResponseData `json:"data" binding:"required"`
	Paging Pagination                       `json:"paging" binding:"required"`
}
