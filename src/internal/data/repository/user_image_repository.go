package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type UserImageRepository struct {
	DBConn *gorm.DB
}

// Create a new record
func (r *UserImageRepository) Create(userImage *entity.SUserImage) error {
	return r.DBConn.Create(userImage).Error
}

// Find by ID
func (r *UserImageRepository) FindByID(id int) (*entity.SUserImage, error) {
	var userImage entity.SUserImage
	if err := r.DBConn.First(&userImage, id).Error; err != nil {
		return nil, err
	}
	return &userImage, nil
}

// Find all
func (r *UserImageRepository) FindAll() ([]entity.SUserImage, error) {
	var userImages []entity.SUserImage
	if err := r.DBConn.Find(&userImages).Error; err != nil {
		return nil, err
	}
	return userImages, nil
}

// Update
func (r *UserImageRepository) Update(userImage *entity.SUserImage) error {
	return r.DBConn.Save(userImage).Error
}

// Delete
func (r *UserImageRepository) Delete(id int) error {
	return r.DBConn.Delete(&entity.SUserImage{}, id).Error
}
