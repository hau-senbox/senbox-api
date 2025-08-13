package request

import "sen-global-api/internal/domain/entity"

type UploadLanguagesConfigRequest struct {
	OwnerID    string                    `json:"owner_id" binding:"required"`
	OwnerRole  string                    `json:"owner_role" binding:"required"`
	SpokenLang entity.LanguageConfigList `json:"spoken_lang" binding:"required"`
	StudyLang  entity.LanguageConfigList `json:"study_lang" binding:"required"`
}
