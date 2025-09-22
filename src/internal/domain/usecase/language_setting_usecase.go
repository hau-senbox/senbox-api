package usecase

import (
	"fmt"
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/response"
)

type LanguageSettingUseCase struct {
	LanguageSettingRepository *repository.LanguageSettingRepository
	ComponentRepository       *repository.ComponentRepository
}

func (l *LanguageSettingUseCase) GetAll() ([]response.LanguageSettingResponse, error) {
	data, err := l.LanguageSettingRepository.GetAll()
	if err != nil {
		return nil, err
	}
	return mapper.MapLanguageSettingToResponses(data), nil
}

func (l *LanguageSettingUseCase) GetAllIsPublished() ([]response.LanguageSettingResponse, error) {
	data, err := l.LanguageSettingRepository.GetAllIsPublished()
	if err != nil {
		return nil, err
	}
	return mapper.MapLanguageSettingToResponses(data), nil
}

func (l *LanguageSettingUseCase) Upload(req request.UploadLanguageSettingRequest) error {
	tx := l.LanguageSettingRepository.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 1. Xóa theo DeleteIDs
	if len(req.DeleteIDs) > 0 {
		if err := l.LanguageSettingRepository.DeleteByIDs(tx, req.DeleteIDs, l.ComponentRepository); err != nil {
			tx.Rollback()
			return err
		}
	}

	// 2. Xử lý create/update LanguageSettings
	for _, s := range req.LanguageSettings {
		if s.ID != nil {
			// Update
			setting, err := l.LanguageSettingRepository.GetByID(*s.ID)
			if err != nil {
				tx.Rollback()
				return err
			}
			// kiem tra exist trong component (neu co) -> khong cho update
			exist, err := l.ComponentRepository.CheckExistLanguage(tx, *s.ID)
			if err != nil {
				tx.Rollback()
				return err
			}
			if exist {
				if setting.LangKey != s.LangKey || setting.RegionKey != s.RegionKey {
					tx.Rollback()
					return fmt.Errorf("language setting ID %d is in use by components and cannot be updated", *s.ID)
				}
			}

			setting.LangKey = s.LangKey
			setting.RegionKey = s.RegionKey
			setting.IsPublished = s.IsPublished

			if err := tx.Save(setting).Error; err != nil {
				tx.Rollback()
				return err
			}
		} else {
			// Create
			setting := &entity.LanguageSetting{
				LangKey:     s.LangKey,
				RegionKey:   s.RegionKey,
				IsPublished: s.IsPublished,
			}
			if err := tx.Create(setting).Error; err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	// 3. Commit transaction
	return tx.Commit().Error
}
