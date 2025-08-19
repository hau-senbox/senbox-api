package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type UserImagesRepository struct {
	DBConn *gorm.DB
}

// Create a new record
func (r *UserImagesRepository) Create(userImage *entity.UserImages) error {
	return r.DBConn.Create(userImage).Error
}

// Get by ID
func (r *UserImagesRepository) GetByID(id string) (*entity.UserImages, error) {
	var userImage entity.UserImages
	if err := r.DBConn.First(&userImage, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &userImage, nil
}

// Get all
func (r *UserImagesRepository) GetAll() ([]entity.UserImages, error) {
	var userImages []entity.UserImages
	if err := r.DBConn.Find(&userImages).Error; err != nil {
		return nil, err
	}
	return userImages, nil
}

// Update
func (r *UserImagesRepository) Update(userImage *entity.UserImages) error {
	return r.DBConn.Save(userImage).Error
}

// Delete
func (r *UserImagesRepository) Delete(id string) error {
	return r.DBConn.Delete(&entity.UserImages{}, "id = ?", id).Error
}

// GetByOwnerAndIndex lấy 1 user image theo owner_id, owner_role và index
func (r *UserImagesRepository) GetByOwnerAndIndex(ownerID string, ownerRole string, index int) (*entity.UserImages, error) {
	var userImage entity.UserImages
	if err := r.DBConn.
		Where("owner_id = ? AND owner_role = ? AND `index` = ?", ownerID, ownerRole, index).
		First(&userImage).Error; err != nil {
		return nil, err
	}
	return &userImage, nil
}

// GetByOwnerAndRole lấy tất cả user_images theo owner_id và owner_role
func (r *UserImagesRepository) GetByOwnerAndRole(ownerID string, ownerRole string) ([]entity.UserImages, error) {
	var userImages []entity.UserImages
	if err := r.DBConn.
		Where("owner_id = ? AND owner_role = ?", ownerID, ownerRole).
		Find(&userImages).Error; err != nil {
		return nil, err
	}
	return userImages, nil
}

func (r *UserImagesRepository) GetByOwnerRoleAndIndex(ownerID string, ownerRole string, index int) (*entity.UserImages, error) {
	var userImage entity.UserImages
	err := r.DBConn.Where("owner_id = ? AND owner_role = ? AND `index` = ?", ownerID, ownerRole, index).
		First(&userImage).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userImage, nil
}

func (r *UserImagesRepository) ResetIsMain(ownerID string, ownerRole string) error {
	err := r.DBConn.Model(&entity.UserImages{}).
		Where("owner_id = ? AND owner_role = ?", ownerID, ownerRole).
		Update("is_main", false).Error
	if err != nil {
		return err
	}
	return nil
}

func (r *UserImagesRepository) GetByOwnerAndRoleIsMain(ownerID string, ownerRole string) (*entity.UserImages, error) {
	var userImage entity.UserImages
	if err := r.DBConn.
		Where("owner_id = ? AND owner_role = ? AND is_main = ?", ownerID, ownerRole, true).
		First(&userImage).Error; err != nil {
		return nil, err
	}
	return &userImage, nil
}
