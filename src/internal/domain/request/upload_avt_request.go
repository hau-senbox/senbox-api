package request

import (
	"sen-global-api/internal/domain/value"
)

type UploadAvatarRequest struct {
	OwnerID   string          `json:"owner_id"`
	OwnerRole value.OwnerRole `json:"owner_role"`
	ImageID   uint64          `json:"image_id"`
	Index     int             `json:"index"`
}
