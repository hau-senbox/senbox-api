package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
)

func ToUserBlockSettingResponse(e *entity.UserBlockSetting) *response.UserBlockSettingResponse {
	return &response.UserBlockSettingResponse{
		ID:              e.ID,
		UserID:          e.UserID,
		IsDeactive:      e.IsDeactive,
		IsViewMessage:   e.IsViewMessage,
		MessageBox:      e.MessageBox,
		MessageDeactive: e.MessageDeactive,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
	}
}
