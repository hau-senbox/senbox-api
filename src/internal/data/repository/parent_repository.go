package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type ParentRepository struct {
	DBConn *gorm.DB
}

func NewParentRepository(db *gorm.DB) *ParentRepository {
	return &ParentRepository{DBConn: db}
}

// Create a new parent record
func (r *ParentRepository) Create(parent *entity.SParent) error {
	return r.DBConn.Create(parent).Error
}

func (r *ParentRepository) GetByUserID(userID string) (*entity.SParent, error) {
	var parents *entity.SParent
	err := r.DBConn.Where("user_id = ?", userID).First(&parents).Error
	return parents, err
}
