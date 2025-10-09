package response

import "sen-global-api/internal/domain/entity"

type LanguagesConfigResponse struct {
	OwnerID    string                  `json:"owner_id"`
	OwnerRole  string                  `json:"owner_role"`
	SpokenLang []entity.LanguageConfig `json:"spoken_lang"`
	StudyLang  []entity.LanguageConfig `json:"study_lang"`
}

type LanguagesConfigResponse4Web struct {
	StudyLang []StudyLanguageConfig4Web `json:"study_lang"`
}

type LanguagesConfigResponse4App struct {
	StudyLang []StudyLanguageConfig4App `json:"study_lang"`
}

type StudyLanguageConfig4Web struct {
	LanguageKey string `json:"lang_key"`
	RegionKey   string `json:"region_key"`
}

type StudyLanguageConfig4App struct {
	Language string `json:"language"`
	Origin   string `json:"origin"`
}
