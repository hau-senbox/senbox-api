package mapper

import (
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/response"
	"strings"
	"unicode"
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
			Order:       item.Order,
			LanguageKey: item.LanguageKey,
			RegionKey:   item.RegionKey,
			Percent:     item.Percent,
			Note:        item.Note,
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
			LanguageKey:       item.LanguageKey,
			RegionKey:         item.RegionKey,
			UniqueLanguageKey: item.LanguageKey + "-" + item.RegionKey,
		}
	}
	return result
}

func ToLanguagesConfigResponse4App(lc *entity.LanguagesConfig) *response.LanguagesConfigResponse4App {
	if lc == nil {
		return nil
	}

	return &response.LanguagesConfigResponse4App{
		StudyLang: toLangItemList4App(lc.StudyLang),
	}
}

func toLangItemList4App(items []entity.LanguageConfig) []response.StudyLanguageConfig4App {
	result := make([]response.StudyLanguageConfig4App, len(items))
	for i, item := range items {
		result[i] = response.StudyLanguageConfig4App{
			Language:          item.LanguageKey,
			Origin:            formatOrigin(item.RegionKey),
			UniqueLanguageKey: item.LanguageKey + "-" + item.RegionKey,
		}
	}
	return result
}

func formatOrigin(origin string) string {
	parts := strings.Split(origin, "_")
	for i, p := range parts {
		if len(p) == 0 {
			continue
		}
		// Viết hoa chữ cái đầu tiên
		runes := []rune(p)
		runes[0] = unicode.ToUpper(runes[0])
		parts[i] = string(runes)
	}
	return strings.Join(parts, " ")
}
