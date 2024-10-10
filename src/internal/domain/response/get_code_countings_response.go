package response

import "sen-global-api/internal/domain/entity"

type GetCodeCountingsResponse struct {
	Codes  []entity.SCodeCounting `json:"codes" binding:"required"`
	Paging Pagination             `json:"pagination" binding:"required"`
}
