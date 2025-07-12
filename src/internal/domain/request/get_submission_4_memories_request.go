package request

type GetSubmission4MemmoriesRequest struct {
	QrCode string `json:"qr_code" binding:"required"`
}
