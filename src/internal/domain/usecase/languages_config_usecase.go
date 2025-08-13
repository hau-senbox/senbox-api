package usecase

import (
	"context"
	"errors"
	"time"

	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/google/uuid"
)

type LanguagesConfigUsecase struct {
	repo *repository.LanguagesConfigRepository
}

func NewLanguagesConfigUsecase(repo *repository.LanguagesConfigRepository) *LanguagesConfigUsecase {
	return &LanguagesConfigUsecase{repo: repo}
}

// Upsert
func (uc *LanguagesConfigUsecase) UploadLanguagesConfig(
	ctx context.Context,
	ownerID string, ownerRole value.OwnerRole4LangConfig,
	spokenLang entity.LanguageConfigList,
	studyLang entity.LanguageConfigList,
) error {

	if ownerID == "" || ownerRole == "" {
		return errors.New("owner id and Owner role is required")
	}

	if len(spokenLang) == 0 {
		return errors.New("spoken Language not empty")
	}

	if len(studyLang) == 0 {
		return errors.New("study Language not empty")
	}

	// Check tồn tại
	existing, err := uc.repo.GetByOwner(ctx, ownerID, ownerRole)
	if err != nil {
		return err
	}

	if existing != nil {
		existing.SpokenLang = spokenLang
		existing.StudyLang = studyLang
		existing.UpdatedAt = time.Now()

		if err := uc.repo.Update(ctx, existing); err != nil {
			return err
		}
		return nil
	}

	// Chưa tồn tại → tạo mới
	newLC := &entity.LanguagesConfig{
		ID:         uuid.New(),
		OwnerID:    ownerID,
		OwnerRole:  ownerRole,
		SpokenLang: spokenLang,
		StudyLang:  studyLang,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := uc.repo.Create(ctx, newLC); err != nil {
		return err
	}
	return nil
}

// GetLanguagesConfigByID - Lấy theo ID
func (uc *LanguagesConfigUsecase) GetLanguagesConfigByID(ctx context.Context, id string) (*entity.LanguagesConfig, error) {
	return uc.repo.GetByID(ctx, id)
}

// GetLanguagesConfigByOwner - Lấy theo OwnerID & OwnerType
func (uc *LanguagesConfigUsecase) GetLanguagesConfigByOwner(ctx context.Context, ownerID string, ownerRole value.OwnerRole4LangConfig) (*response.LanguagesConfigResponse, error) {

	lc, err := uc.repo.GetByOwner(ctx, ownerID, ownerRole)
	if err != nil {
		return nil, err
	}
	return mapper.ToLanguagesConfigResponse(lc), nil
}

// UpdateLanguagesConfig - Cập nhật
func (uc *LanguagesConfigUsecase) UpdateLanguagesConfig(ctx context.Context, lc *entity.LanguagesConfig) error {
	if lc.ID == uuid.Nil {
		return errors.New("ID không được rỗng khi update")
	}
	lc.UpdatedAt = time.Now()
	return uc.repo.Update(ctx, lc)
}

// DeleteLanguagesConfig - Xoá
func (uc *LanguagesConfigUsecase) DeleteLanguagesConfig(ctx context.Context, id string) error {
	return uc.repo.Delete(ctx, id)
}
