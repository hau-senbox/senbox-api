package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/value"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type StaffApplicationRepository struct {
	DBConn *gorm.DB
}

func NewStaffApplicationRepository(db *gorm.DB) *StaffApplicationRepository {
	return &StaffApplicationRepository{DBConn: db}
}

// Create a new staff application
func (r *StaffApplicationRepository) Create(app *entity.SStaffFormApplication) error {
	return r.DBConn.Create(app).Error
}

// Get by ID
func (r *StaffApplicationRepository) GetByID(id uuid.UUID) (*entity.SStaffFormApplication, error) {
	var app entity.SStaffFormApplication

	err := r.DBConn.Where("id = ?", id).First(&app).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// Get all applications
func (r *StaffApplicationRepository) GetAll() ([]entity.SStaffFormApplication, error) {
	var apps []entity.SStaffFormApplication
	err := r.DBConn.Find(&apps).Error
	return apps, err
}

// Get all applications approved
func (r *StaffApplicationRepository) GetApprovedAll() ([]entity.SStaffFormApplication, error) {
	var apps []entity.SStaffFormApplication

	err := r.DBConn.
		Where("status = ?", value.Approved).
		Find(&apps).Error

	return apps, err
}

// Update application
func (r *StaffApplicationRepository) Update(app *entity.SStaffFormApplication) error {
	return r.DBConn.Save(app).Error
}

// Delete by ID
func (r *StaffApplicationRepository) Delete(id int64) error {
	return r.DBConn.Delete(&entity.SStaffFormApplication{}, id).Error
}

// Get by UserID
func (r *StaffApplicationRepository) GetByUserIDApproved(userID string) ([]entity.SStaffFormApplication, error) {
	var apps []entity.SStaffFormApplication
	err := r.DBConn.Where("user_id = ? AND status = ?", userID, value.Approved).Find(&apps).Error
	return apps, err
}

// Get by OrganizationID
func (r *StaffApplicationRepository) GetByOrganizationID(orgID string) ([]entity.SStaffFormApplication, error) {
	var apps []entity.SStaffFormApplication
	err := r.DBConn.Where("organization_id = ? AND status = ?", orgID, value.Approved).Find(&apps).Error
	return apps, err
}

// GetAllStaffIDs returns a list of all staff application IDs
func (r *StaffApplicationRepository) GetAllStaffIDs() ([]uuid.UUID, error) {
	var ids []uuid.UUID

	err := r.DBConn.
		Model(&entity.SStaffFormApplication{}).
		Pluck("id", &ids).Error

	if err != nil {
		return nil, err
	}

	return ids, nil
}

func (r *StaffApplicationRepository) GetDB() *gorm.DB {
	return r.DBConn
}

func (r *StaffApplicationRepository) GetByOrganizationIDsApproved(orgIDStrs []string) ([]entity.SStaffFormApplication, error) {
	var apps []entity.SStaffFormApplication

	if len(orgIDStrs) == 0 {
		return []entity.SStaffFormApplication{}, nil
	}

	err := r.DBConn.Where("organization_id IN ? AND status = ?", orgIDStrs, value.Approved).Find(&apps).Error
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *StaffApplicationRepository) GetByOrganizationIDs(orgIDStrs []string) ([]entity.SStaffFormApplication, error) {
	var apps []entity.SStaffFormApplication

	if len(orgIDStrs) == 0 {
		return []entity.SStaffFormApplication{}, nil
	}

	err := r.DBConn.Where("organization_id IN ?", orgIDStrs).Find(&apps).Error
	if err != nil {
		return nil, err
	}

	return apps, nil
}

func (r *StaffApplicationRepository) CheckStaffBelongsToOrganizations(tx *gorm.DB, staffID uuid.UUID, orgIDs []string) (bool, error) {
	var count int64
	err := tx.Model(&entity.SStaffFormApplication{}).
		Where("id = ? AND organization_id IN ?", staffID, orgIDs).
		Count(&count).Error
	return count > 0, err
}

func (r *StaffApplicationRepository) GetByUserID(userID string) (entity.SStaffFormApplication, error) {
	var app entity.SStaffFormApplication
	err := r.DBConn.Where("user_id = ? AND status = ?", userID, value.Approved).First(&app).Error
	return app, err
}

func (r *StaffApplicationRepository) GetAllByUserID(userID string) ([]entity.SStaffFormApplication, error) {
	var app []entity.SStaffFormApplication
	err := r.DBConn.Where("user_id = ? AND status = ?", userID, value.Approved).Find(&app).Error
	return app, err
}

func (r *StaffApplicationRepository) GetByUserIDAndOrgID(userID string, orgID string) (*entity.SStaffFormApplication, error) {
	var app *entity.SStaffFormApplication
	err := r.DBConn.Where("user_id = ? AND organization_id = ? AND status = ?", userID, orgID, value.Approved).First(&app).Error
	return app, err
}
