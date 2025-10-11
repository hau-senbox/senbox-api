package usecase

import (
	"context"
	"errors"
	"fmt"
	"time"

	"sen-global-api/internal/data/repository"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/mapper"
	"sen-global-api/internal/domain/response"
	"sen-global-api/internal/domain/value"

	"github.com/google/uuid"
)

type LanguagesConfigUsecase struct {
	Repo            *repository.LanguagesConfigRepository
	LangSettingRepo *repository.LanguageSettingRepository
}

// func NewLanguagesConfigUsecase(repo *repository.LanguagesConfigRepository) *LanguagesConfigUsecase {
// 	return &LanguagesConfigUsecase{Repo: repo}
// }

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

	// Validate tổng percent
	if sum := spokenLang.TotalPercent(); sum > 100 {
		return fmt.Errorf("spoken languages percent exceeded 100 (got %d)", sum)
	}
	if sum := studyLang.TotalPercent(); sum > 100 {
		return fmt.Errorf("study languages percent exceeded 100 (got %d)", sum)
	}

	// Check tồn tại
	existing, err := uc.Repo.GetByOwner(ctx, ownerID, ownerRole)
	if err != nil {
		return err
	}

	if existing != nil {
		existing.SpokenLang = spokenLang
		existing.StudyLang = studyLang
		existing.UpdatedAt = time.Now()

		if err := uc.Repo.Update(ctx, existing); err != nil {
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
	if err := uc.Repo.Create(ctx, newLC); err != nil {
		return err
	}
	return nil
}

// GetLanguagesConfigByID - Lấy theo ID
func (uc *LanguagesConfigUsecase) GetLanguagesConfigByID(ctx context.Context, id string) (*entity.LanguagesConfig, error) {
	return uc.Repo.GetByID(ctx, id)
}

// GetLanguagesConfigByOwner - Lấy theo OwnerID & OwnerType
func (uc *LanguagesConfigUsecase) GetLanguagesConfigByOwner(ctx context.Context, ownerID string, ownerRole value.OwnerRole4LangConfig) (*response.LanguagesConfigResponse, error) {

	lc, err := uc.Repo.GetByOwner(ctx, ownerID, ownerRole)
	if err != nil {
		return nil, err
	}
	return mapper.ToLanguagesConfigResponse(lc), nil
}

func (uc *LanguagesConfigUsecase) GetLanguagesConfigByOwnerNoCtx(ownerID string, ownerRole value.OwnerRole4LangConfig) (*response.LanguagesConfigResponse, error) {

	lc, err := uc.Repo.GetByOwnerNoCtx(ownerID, ownerRole)
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
	return uc.Repo.Update(ctx, lc)
}

// DeleteLanguagesConfig - Xoá
func (uc *LanguagesConfigUsecase) DeleteLanguagesConfig(ctx context.Context, id string) error {
	return uc.Repo.Delete(ctx, id)
}

func (uc *LanguagesConfigUsecase) GetStudentStudyLangConfig(ctx context.Context, studentID string) (*response.LanguagesConfigResponse, error) {

	lc, err := uc.Repo.GetByOwner(ctx, studentID, value.OwnerRoleLangStudent)
	if err != nil {
		return nil, err
	}
	return mapper.ToStudyLanguagesConfigResponse(lc), nil
}

func (uc *LanguagesConfigUsecase) GetLanguagesConfigByOwner4Web(ctx context.Context, ownerID string, ownerRole value.OwnerRole4LangConfig) (*response.LanguagesConfigResponse4Web, error) {

	lc, err := uc.Repo.GetByOwner(ctx, ownerID, ownerRole)
	if err != nil {
		return nil, err
	}
	return mapper.ToLanguagesConfigResponse4Web(lc), nil
}

func (uc *LanguagesConfigUsecase) GetLanguagesConfigByOwner4App(ctx context.Context, ownerID string, ownerRole value.OwnerRole4LangConfig) (*response.LanguagesConfigResponse4App, error) {

	lc, err := uc.Repo.GetByOwner(ctx, ownerID, ownerRole)
	if err != nil {
		return nil, err
	}
	return mapper.ToLanguagesConfigResponse4App(lc), nil
}

func (uc *LanguagesConfigUsecase) GetStudyLanguage4OrganizationAssign4Web(ctx context.Context, organizationID string) ([]*response.StudyLanguage4Assign, error) {

	// get list language setting form admin
	languageSetting, err := uc.LangSettingRepo.GetAllIsPublished()
	if err != nil {
		return nil, errors.New("cannot get language setting")
	}

	// get list language config form organization
	languageConfig, err := uc.Repo.GetByOwner(ctx, organizationID, value.OwnerRoleLangOrganization)

	if err != nil {
		return nil, errors.New("cannot get language config")
	}

	return mapper.ToAssignLanguagesConfigResponse(languageSetting, &languageConfig.StudyLang), nil

}
