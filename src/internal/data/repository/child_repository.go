package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChildRepository struct {
	DB *gorm.DB
}

func NewChildRepository(db *gorm.DB) *ChildRepository {
	return &ChildRepository{DB: db}
}

// Create a new child
func (r *ChildRepository) Create(child *entity.SChild) error {
	var existing entity.SChild

	// Kiểm tra tên trùng trong cùng 1 user
	err := r.DB.
		Where("child_name = ? AND parent_id = ?", child.ChildName, child.ParentID).
		First(&existing).Error

	if err == nil {
		return errors.New("child with the same name already exists for this user")
	}

	// Nếu lỗi không phải là "record not found" => báo lỗi hệ thống
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return err
	}

	// Nếu chưa có thì tạo mới
	child.ID = uuid.New()
	return r.DB.Create(child).Error
}

// Get all children by parent ID
func (r *ChildRepository) GetByParentID(userID string) ([]entity.SChild, error) {
	var children []entity.SChild
	err := r.DB.Where("parent_id = ?", userID).Find(&children).Error
	return children, err
}

// Get child by ID
func (r *ChildRepository) GetByID(id string) (*entity.SChild, error) {
	var child entity.SChild

	parsedID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("invalid UUID format")
	}

	err = r.DB.Where("id = ?", parsedID).First(&child).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &child, nil
}

// Update child
func (r *ChildRepository) Update(child *entity.SChild) error {
	return r.DB.Save(child).Error
}

// Delete child
func (r *ChildRepository) Delete(id int64) error {
	return r.DB.Delete(&entity.SChild{}, id).Error
}
