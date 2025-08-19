package repository

import (
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
