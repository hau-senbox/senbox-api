package repository

import (
	"errors"
	"fmt"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
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

// GetAllIDs returns a list of all child IDs
func (r *ChildRepository) GetAllIDs() ([]uuid.UUID, error) {
	var ids []uuid.UUID
	err := r.DB.Model(&entity.SChild{}).Select("id").Find(&ids).Error
	return ids, err
}

// GetAll returns all children
func (r *ChildRepository) GetAll() ([]entity.SChild, error) {
	var children []entity.SChild
	err := r.DB.Find(&children).Error
	return children, err
}

func (r *ChildMenuRepository) DeleteAllTx(tx *gorm.DB) error {
	if err := tx.Exec("DELETE FROM child_menu").Error; err != nil {
		log.Error("ChildMenuRepository.DeleteAllTx: " + err.Error())
		return fmt.Errorf("Delete all child_menu fail: %w", err)
	}
	return nil
}

func (r *ChildRepository) GetParentIDByChildID(childID string) (string, error) {
	var child entity.SChild

	err := r.DB.
		Select("parent_id").
		Where("id = ?", childID).
		First(&child).Error

	if err != nil {
		return "", err
	}

	return child.ParentID.String(), nil
}

func (r *ChildRepository) GetAllParents() ([]entity.SUserEntity, error) {
	var parents []entity.SUserEntity

	err := r.DB.
		Table("s_child AS c").
		Select("DISTINCT u.*").
		Joins(`
			JOIN s_user_entity AS u 
			ON CONVERT(u.id USING utf8mb4) COLLATE utf8mb4_unicode_ci = 
			   CONVERT(c.parent_id USING utf8mb4) COLLATE utf8mb4_unicode_ci
		`).
		Scan(&parents).Error

	if err != nil {
		return nil, err
	}

	return parents, nil
}
