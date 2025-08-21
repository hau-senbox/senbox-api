package request

type UploadAppConfigRequest struct {
	Key   string `json:"key" binding:"required"`
	Value any    `json:"value" binding:"required"`
}
