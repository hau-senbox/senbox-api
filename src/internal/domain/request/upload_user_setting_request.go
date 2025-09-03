package request

type UploadUserSettingRequest struct {
	OwnerID   string `json:"owner_id" binding:"required"`
	OwnerRole string `json:"owner_role" binding:"required"`
	Key       string `json:"key" binding:"required"`
	Value     any    `json:"value" binding:"required"`
}
