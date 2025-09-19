package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"
)

// Map 1 UserSetting entity -> UserSettingResponse
func ToUserSettingResponse(list []*entity.UserSetting) *response.UserSettingResponse {
	if list == nil {
		return nil
	}
	res := &response.UserSettingResponse{}
	for _, e := range list {
		switch e.Key {
		case value.UserSettingLoginDeviceLimit:
			res.LimitDeviceLogin = e.Value
		case value.UserSettingLanguage:
			res.AppLanguage = e.Value
		}
	}
	return res
}
