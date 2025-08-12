package mapper

import (
	"fmt"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/response"
	"time"
)

func MapOrgSettingToResponse(setting *entity.OrganizationSetting, comp *components.Component) response.OrgSettingResponse {

	return response.OrgSettingResponse{
		ID:                fmt.Sprintf("%d", setting.ID),
		OrganizationID:    setting.OrganizationID,
		IsViewMessage:     setting.IsViewMessage,
		IsShowOrgNews:     setting.IsShowOrgNews,
		IsDeactiveTopMenu: setting.IsDeactiveTopMenu,
		IsShowSpecialBtn:  setting.IsShowSpecialBtn,
		MessageBox:        setting.MessageBox,
		MessageTopMenu:    setting.MessageTopMenu,
		Component:         comp,
		CreatedAt:         setting.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         setting.UpdatedAt.Format(time.RFC3339),
	}
}
