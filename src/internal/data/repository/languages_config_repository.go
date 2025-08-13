package repository

import (
	"context"
	"errors"

	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"

	"gorm.io/gorm"
)

type LanguagesConfigRepository struct {
	db *gorm.DB
}

func NewLanguagesConfigRepository(db *gorm.DB) *LanguagesConfigRepository {
	return &LanguagesConfigRepository{db: db}
}

func (r *LanguagesConfigRepository) Create(ctx context.Context, lc *entity.LanguagesConfig) error {
	return r.db.WithContext(ctx).Create(lc).Error
}

func (r *LanguagesConfigRepository) GetByID(ctx context.Context, id string) (*entity.LanguagesConfig, error) {
	var lc entity.LanguagesConfig
	err := r.db.WithContext(ctx).First(&lc, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // không tìm thấy
		}
		return nil, err
	}
	return &lc, nil
}

func (r *LanguagesConfigRepository) GetByOwner(ctx context.Context, ownerID string, ownerRole value.OwnerRole4LangConfig) (*entity.LanguagesConfig, error) {
	var lc entity.LanguagesConfig
	err := r.db.WithContext(ctx).
		Where("owner_id = ? AND owner_role = ?", ownerID, ownerRole).
		First(&lc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil // không có thì trả về nil
		}
		return nil, err
	}
	return &lc, nil
}

func (r *LanguagesConfigRepository) Update(ctx context.Context, lc *entity.LanguagesConfig) error {
	return r.db.WithContext(ctx).
		Model(&entity.LanguagesConfig{}).
		Where("id = ?", lc.ID).
		Updates(lc).Error
}

func (r *LanguagesConfigRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&entity.LanguagesConfig{}, "id = ?", id).Error
}
