package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type StudentApplicationRepository struct {
	DB *gorm.DB
}

func NewStudentApplicationRepository(db *gorm.DB) *StudentApplicationRepository {
	return &StudentApplicationRepository{DB: db}
}

// Create a new student application
func (r *StudentApplicationRepository) Create(app *entity.SStudentFormApplication) error {
	return r.DB.Create(app).Error
}

// Get by ID
func (r *StudentApplicationRepository) GetByID(id int64) (*entity.SStudentFormApplication, error) {
	var app entity.SStudentFormApplication
	err := r.DB.First(&app, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &app, err
}

// Get all applications
func (r *StudentApplicationRepository) GetAll() ([]entity.SStudentFormApplication, error) {
	var apps []entity.SStudentFormApplication
	err := r.DB.Find(&apps).Error
	return apps, err
}

// Update application
func (r *StudentApplicationRepository) Update(app *entity.SStudentFormApplication) error {
	return r.DB.Save(app).Error
}

// Delete by ID
func (r *StudentApplicationRepository) Delete(id int64) error {
	return r.DB.Delete(&entity.SStudentFormApplication{}, id).Error
}

// Get by UserID
func (r *StudentApplicationRepository) GetByUserID(userID string) ([]entity.SStudentFormApplication, error) {
	var apps []entity.SStudentFormApplication
	err := r.DB.Where("user_id = ?", userID).Find(&apps).Error
	return apps, err
}

// Get by OrganizationID
func (r *StudentApplicationRepository) GetByOrganizationID(orgID string) ([]entity.SStudentFormApplication, error) {
	var apps []entity.SStudentFormApplication
	err := r.DB.Where("organization_id = ?", orgID).Find(&apps).Error
	return apps, err
}
