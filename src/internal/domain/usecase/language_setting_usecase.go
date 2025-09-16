package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type LanguageSettingUseCase struct {
	LanguageSettingRepository *repository.LanguageSettingRepository
}

func (l *LanguageSettingUseCase) GetAll() ([]entity.LanguageSetting, error) {
	return l.LanguageSettingRepository.GetAll()
}

func (l *LanguageSettingUseCase) Upload(req request.UploadLanguageSettingRequest) error {
	return nil
}
