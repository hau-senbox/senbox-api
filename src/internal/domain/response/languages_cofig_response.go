package response

import "sen-global-api/internal/domain/entity"

type LanguagesConfigResponse struct {
	OwnerID    string                  `json:"owner_id"`
	OwnerRole  string                  `json:"owner_role"`
	SpokenLang []entity.LanguageConfig `json:"spoken_lang"`
	StudyLang  []entity.LanguageConfig `json:"study_lang"`
}
