package repository

import (
	"context"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type ParentChildsRepository struct {
	DBConn *gorm.DB
}

func NewParentChildsRepository(DBConn *gorm.DB) ParentChildsRepository {
	return ParentChildsRepository{DBConn: DBConn}
}

func (r *ParentChildsRepository) Create(ctx context.Context, rel *entity.SParentChilds) error {
	return r.DBConn.WithContext(ctx).Create(rel).Error
}

func (r *ParentChildsRepository) GetByParentID(ctx context.Context, parentID string) ([]entity.SParentChilds, error) {
	var result []entity.SParentChilds
	err := r.DBConn.WithContext(ctx).Where("parent_id = ?", parentID).Find(&result).Error
	return result, err
}

func (r *ParentChildsRepository) GetByChildID(ctx context.Context, childID string) ([]entity.SParentChilds, error) {
	var result []entity.SParentChilds
	err := r.DBConn.WithContext(ctx).Where("child_id = ?", childID).Find(&result).Error
	return result, err
}

func (r *ParentChildsRepository) Delete(ctx context.Context, parentID, childID string) error {
	return r.DBConn.WithContext(ctx).
		Where("parent_id = ? AND child_id = ?", parentID, childID).
		Delete(&entity.SParentChilds{}).Error
}
