package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/entity/components"
	"sen-global-api/internal/domain/response"
)

func MapOrgSettingToResponse(setting *entity.OrganizationSetting, comp *components.Component) response.OrgSettingResponse {

	return response.OrgSettingResponse{
		ID:                 setting.ID.String(),
		OrganizationID:     setting.OrganizationID,
		DeviceID:           setting.DeviceID,
		IsViewMessageBox:   setting.IsViewMessageBox,
		IsShowMessage:      setting.IsShowMessage,
		MessageBox:         setting.MessageBox,
		IsShowSpecialBtn:   setting.IsShowSpecialBtn,
		IsDeactiveApp:      setting.IsDeactiveApp,
		MessageDeactiveApp: setting.MessageDeactiveApp,
		IsDeactiveTopMenu:  setting.IsDeactiveTopMenu,
		MessageTopMenu:     setting.MessageTopMenu,
		TopMenuPassword:    setting.TopMenuPassword,
		Component:          comp,
	}
}
