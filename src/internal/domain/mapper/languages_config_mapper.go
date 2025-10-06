package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
)

func ToLanguagesConfigResponse(lc *entity.LanguagesConfig) *response.LanguagesConfigResponse {
	if lc == nil {
		return nil
	}

	return &response.LanguagesConfigResponse{
		OwnerID:    lc.OwnerID,
		OwnerRole:  string(lc.OwnerRole),
		SpokenLang: toLangItemList(lc.SpokenLang),
		StudyLang:  toLangItemList(lc.StudyLang),
	}
}

func ToStudyLanguagesConfigResponse(lc *entity.LanguagesConfig) *response.LanguagesConfigResponse {
	if lc == nil {
		return nil
	}

	return &response.LanguagesConfigResponse{
		OwnerID:   lc.OwnerID,
		OwnerRole: string(lc.OwnerRole),
		StudyLang: toLangItemList(lc.StudyLang),
	}
}

func toLangItemList(items []entity.LanguageConfig) []entity.LanguageConfig {
	result := make([]entity.LanguageConfig, len(items))
	for i, item := range items {
		result[i] = entity.LanguageConfig{
			Order:    item.Order,
			Language: item.Language,
			Origin:   item.Origin,
			Percent:  item.Percent,
			Note:     item.Note,
		}
	}
	return result
}

func ToLanguagesConfigResponse4Web(lc *entity.LanguagesConfig) *response.LanguagesConfigResponse4Web {
	if lc == nil {
		return nil
	}

	return &response.LanguagesConfigResponse4Web{
		StudyLang: toLangItemList4Web(lc.StudyLang),
	}
}

func toLangItemList4Web(items []entity.LanguageConfig) []response.StudyLanguageConfig4Web {
	result := make([]response.StudyLanguageConfig4Web, len(items))
	for i, item := range items {
		result[i] = response.StudyLanguageConfig4Web{
			LanguageKey: item.Language,
			RegionKey:   item.Origin,
		}
	}
	return result
}
