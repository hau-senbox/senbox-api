package repository

import (
	"context"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ParentRepository struct {
	DBConn *gorm.DB
}

func NewParentRepository(db *gorm.DB) *ParentRepository {
	return &ParentRepository{DBConn: db}
}

// Create a new parent record
func (r *ParentRepository) Create(ctx context.Context, parent *entity.SParent) error {
	return r.DBConn.WithContext(ctx).Create(parent).Error
}

func (r *ParentRepository) GetByUserID(ctx context.Context, userID string) (*entity.SParent, error) {
	var parents *entity.SParent
	err := r.DBConn.WithContext(ctx).Where("user_id = ?", userID).First(&parents).Error
	return parents, err
}

func (r *ParentRepository) WithTx(tx *gorm.DB) *ParentRepository {
	return &ParentRepository{DBConn: tx}
}

func (r *ParentRepository) GetAll(ctx context.Context) ([]entity.SParent, error) {
	var parents []entity.SParent
	err := r.DBConn.WithContext(ctx).Find(&parents).Error
	return parents, err
}

func (r *ParentRepository) GetByID(ctx context.Context, parentID string) (*entity.SParent, error) {
	var parents *entity.SParent

	parentUuid := uuid.MustParse(parentID)
	err := r.DBConn.WithContext(ctx).Where("id = ?", parentUuid).First(&parents).Error
	return parents, err
}
