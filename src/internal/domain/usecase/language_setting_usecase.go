package usecase

import (
	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"strconv"
)

type LanguageSettingUseCase struct {
	LanguageSettingRepository *repository.LanguageSettingRepository
	ComponentRepository       *repository.ComponentRepository
}

func (l *LanguageSettingUseCase) GetAll() ([]entity.LanguageSetting, error) {
	return l.LanguageSettingRepository.GetAll()
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
		var ids []uint
		for _, idStr := range req.DeleteIDs {
			id, err := strconv.ParseUint(idStr, 10, 64)
			if err != nil {
				tx.Rollback()
				return err
			}
			ids = append(ids, uint(id))
		}

		if err := l.LanguageSettingRepository.DeleteByIDs(tx, ids, l.ComponentRepository); err != nil {
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
