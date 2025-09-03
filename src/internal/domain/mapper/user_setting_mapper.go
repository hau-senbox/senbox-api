package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
)

// Map 1 UserSetting entity -> UserSettingResponse
func ToUserSettingResponse(e *entity.UserSetting) *response.UserSettingResponse {
	if e == nil {
		return nil
	}
	return &response.UserSettingResponse{
		Key:   string(e.Key),
		Value: e.Value,
	}
}

// Map list UserSettings -> list UserSettingResponse
func ToUserSettingResponses(list []*entity.UserSetting) []*response.UserSettingResponse {
	if list == nil {
		return nil
	}
	res := make([]*response.UserSettingResponse, 0, len(list))
	for _, e := range list {
		res = append(res, ToUserSettingResponse(e))
	}
	return res
}
