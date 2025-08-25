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
