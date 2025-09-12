package response

import "sen-global-api/internal/domain/value"

type OwnerAssignResponse struct {
	OwnerID   string          `json:"owner_id"`
	OwnerRole value.OwnerRole `json:"owner_role"`
	Name      string          `json:"name"`
	AvatarKey string          `json:"avatar_key"`
	AvatarUrl string          `json:"avatar_url"`
}

type ListOwnerAssignResponse struct {
	Teachers []*OwnerAssignResponse `json:"teachers"`
	Staffs   []*OwnerAssignResponse `json:"staffs"`
	Students []*OwnerAssignResponse `json:"students"`
}
