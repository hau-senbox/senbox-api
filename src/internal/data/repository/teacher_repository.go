package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeacherApplicationRepository struct {
	DBconn *gorm.DB
}

func NewTeacherApplicationRepository(db *gorm.DB) *TeacherApplicationRepository {
	return &TeacherApplicationRepository{DBconn: db}
}

// Create a new teacher application
func (r *TeacherApplicationRepository) Create(app *entity.STeacherFormApplication) error {
	return r.DBconn.Create(app).Error
}

// Get by ID
func (r *TeacherApplicationRepository) GetByID(id uuid.UUID) (*entity.STeacherFormApplication, error) {
	var app entity.STeacherFormApplication
	err := r.DBconn.Where("id = ?", id).First(&app).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Get all applications
func (r *TeacherApplicationRepository) GetAll() ([]entity.STeacherFormApplication, error) {
	var apps []entity.STeacherFormApplication
	err := r.DBconn.Find(&apps).Error
	return apps, err
}

// Update application
func (r *TeacherApplicationRepository) Update(app *entity.STeacherFormApplication) error {
	return r.DBconn.Save(app).Error
}

// Delete by ID
func (r *TeacherApplicationRepository) Delete(id uuid.UUID) error {
	return r.DBconn.Delete(&entity.STeacherFormApplication{}, id).Error
}

// Get by UserID
func (r *TeacherApplicationRepository) GetByUserID(userID string) ([]entity.STeacherFormApplication, error) {
	var apps []entity.STeacherFormApplication
	err := r.DBconn.Where("user_id = ?", userID).Find(&apps).Error
	return apps, err
}

// Get by OrganizationID
func (r *TeacherApplicationRepository) GetByOrganizationID(orgID string) ([]entity.STeacherFormApplication, error) {
	var apps []entity.STeacherFormApplication
	err := r.DBconn.Where("organization_id = ?", orgID).Find(&apps).Error
	return apps, err
}

// Get by list of OrganizationIDs
func (r *TeacherApplicationRepository) GetByOrganizationIDs(orgIDStrs []string) ([]entity.STeacherFormApplication, error) {
	var apps []entity.STeacherFormApplication
	if len(orgIDStrs) == 0 {
		return []entity.STeacherFormApplication{}, nil
	}

	err := r.DBconn.Where("organization_id IN ?", orgIDStrs).Find(&apps).Error
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// Check if teacher belongs to one of the given organizations
func (r *TeacherApplicationRepository) CheckTeacherBelongsToOrganizations(tx *gorm.DB, teacherID uuid.UUID, orgIDs []string) (bool, error) {
	var count int64
	err := tx.Model(&entity.STeacherFormApplication{}).
		Where("id = ? AND organization_id IN ?", teacherID, orgIDs).
		Count(&count).Error
	return count > 0, err
}

// Get GORM instance (optional)
func (r *TeacherApplicationRepository) GetDB() *gorm.DB {
	return r.DBconn
}

// GetAllTeacherIDs returns a list of all teacher application IDs
func (r *TeacherApplicationRepository) GetAllTeacherIDs() ([]uuid.UUID, error) {
	var ids []uuid.UUID

	err := r.DBconn.
		Model(&entity.STeacherFormApplication{}).
		Pluck("id", &ids).Error

	if err != nil {
		return nil, err
	}

	return ids, nil
}
