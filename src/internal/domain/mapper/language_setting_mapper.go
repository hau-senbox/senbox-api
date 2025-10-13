package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
)

func MapLanguageSettingToResponses(list []entity.LanguageSetting) []response.LanguageSettingResponse {
	result := make([]response.LanguageSettingResponse, len(list))
	for i, item := range list {
		result[i] = response.LanguageSettingResponse{
			ID:                 item.ID,
			LangKey:            item.LangKey,
			RegionKey:          item.RegionKey,
			IsPublished:        item.IsPublished,
			DeactivatedMessage: item.DeactivatedMessage,
			UniqueLanguageKey:  item.LangKey + "-" + item.RegionKey,
		}
	}
	return result
}
