package response

import "sen-global-api/internal/domain/entity/components"

type GetCommonMenuResponse struct {
	Component []components.Component `json:"components"`
}
