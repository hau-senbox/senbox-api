package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TeacherApplicationRepository struct {
	DBConn *gorm.DB
}

func NewTeacherApplicationRepository(db *gorm.DB) *TeacherApplicationRepository {
	return &TeacherApplicationRepository{DBConn: db}
}

// Create a new teacher application
func (r *TeacherApplicationRepository) Create(app *entity.STeacherFormApplication) error {
	return r.DBConn.Create(app).Error
}

// Get by ID
func (r *TeacherApplicationRepository) GetByID(id uuid.UUID) (*entity.STeacherFormApplication, error) {
	var app entity.STeacherFormApplication
	err := r.DBConn.Where("id = ?", id).First(&app).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, errors.New("teacher not found")
	}
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Get all applications
func (r *TeacherApplicationRepository) GetAll() ([]entity.STeacherFormApplication, error) {
	var apps []entity.STeacherFormApplication
	err := r.DBConn.Find(&apps).Error
	return apps, err
}

func (r *TeacherApplicationRepository) GetApprovedAll() ([]entity.STeacherFormApplication, error) {
	var apps []entity.STeacherFormApplication

	err := r.DBConn.
		Where("status = ?", value.Approved).
		Find(&apps).Error

	return apps, err
}

// Update application
func (r *TeacherApplicationRepository) Update(app *entity.STeacherFormApplication) error {
	return r.DBConn.Save(app).Error
}

// Delete by ID
func (r *TeacherApplicationRepository) Delete(id uuid.UUID) error {
	return r.DBConn.Delete(&entity.STeacherFormApplication{}, id).Error
}

// Get by UserID
func (r *TeacherApplicationRepository) GetByUserIDApproved(userID string) ([]entity.STeacherFormApplication, error) {
	var apps []entity.STeacherFormApplication
	err := r.DBConn.Where("user_id = ? AND status = ?", userID, value.Approved).Find(&apps).Error
	return apps, err
}

// Get by UserID
func (r *TeacherApplicationRepository) GetByUserID(userID string) (entity.STeacherFormApplication, error) {
	var app entity.STeacherFormApplication
	err := r.DBConn.Where("user_id = ? AND status = ?", userID, value.Approved).First(&app).Error
	return app, err
}

// Get all by UserID
func (r *TeacherApplicationRepository) GetAllByUserID(userID string) ([]entity.STeacherFormApplication, error) {
	var app []entity.STeacherFormApplication
	err := r.DBConn.Where("user_id = ? AND status = ?", userID, value.Approved).Find(&app).Error
	return app, err
}

// Get by OrganizationID
func (r *TeacherApplicationRepository) GetByOrganizationID(orgID string) ([]entity.STeacherFormApplication, error) {
	var apps []entity.STeacherFormApplication
	err := r.DBConn.Where("organization_id = ? AND status = ?", orgID, value.Approved).Find(&apps).Error
	return apps, err
}

func (r *TeacherApplicationRepository) GetByOrganizationIDsApproved(orgIDStrs []string) ([]entity.STeacherFormApplication, error) {
	var apps []entity.STeacherFormApplication
	if len(orgIDStrs) == 0 {
		return []entity.STeacherFormApplication{}, nil
	}

	err := r.DBConn.Where("organization_id IN ? AND status = ?", orgIDStrs, value.Approved).Find(&apps).Error
	if err != nil {
		return nil, err
	}

	return apps, nil
}

// Get by list of OrganizationIDs
func (r *TeacherApplicationRepository) GetByOrganizationIDs(orgIDStrs []string) ([]entity.STeacherFormApplication, error) {
	var apps []entity.STeacherFormApplication
	if len(orgIDStrs) == 0 {
		return []entity.STeacherFormApplication{}, nil
	}

	err := r.DBConn.Where("organization_id IN ?", orgIDStrs).Find(&apps).Error
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
	return r.DBConn
}

// GetAllTeacherIDs returns a list of all teacher application IDs
func (r *TeacherApplicationRepository) GetAllTeacherIDs() ([]uuid.UUID, error) {
	var ids []uuid.UUID

	err := r.DBConn.
		Model(&entity.STeacherFormApplication{}).
		Pluck("id", &ids).Error

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *TeacherApplicationRepository) GetByUserIDAndOrgID(userID string, orgID string) (*entity.STeacherFormApplication, error) {
	var app *entity.STeacherFormApplication
	err := r.DBConn.Where("user_id = ? AND organization_id = ? AND status = ?", userID, orgID, value.Approved).First(&app).Error
	return app, err
}
