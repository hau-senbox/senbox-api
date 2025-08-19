package request

import "sen-global-api/internal/domain/value"

type UpdateIsMainAvatar struct {
	OwnerID   string          `json:"owner_id"`
	OwnerRole value.OwnerRole `json:"owner_role"`
	Index     int             `json:"index"`
}
