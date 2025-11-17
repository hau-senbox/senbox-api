package repository

import (
	"context"
	"errors"
	"sen-global-api/internal/domain/entity"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StudentMenuOrganizationRepository struct {
	DBConn *gorm.DB
}

func (r *StudentMenuOrganizationRepository) GetByID(ctx context.Context, id string) (*entity.StudentMenuOrganization, error) {
	var menuOrg entity.StudentMenuOrganization
	if err := r.DBConn.WithContext(ctx).First(&menuOrg, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &menuOrg, nil
}

func (r *StudentMenuOrganizationRepository) GetAllByStudentAndOrg(ctx context.Context, studentID, orgID string) ([]entity.StudentMenuOrganization, error) {
	var menuOrgs []entity.StudentMenuOrganization
	err := r.DBConn.WithContext(ctx).
		Where("student_id = ? AND organization_id = ?", studentID, orgID).
		Order("`order` ASC").
		Find(&menuOrgs).Error
	return menuOrgs, err
}

func (r *StudentMenuOrganizationRepository) GetByStudentOrgAndComponentID(
	tx *gorm.DB,
	studentID, orgID, componentID string,
) (*entity.StudentMenuOrganization, error) {
	var menuOrg entity.StudentMenuOrganization
	if err := tx.Where("student_id = ? AND organization_id = ? AND component_id = ?", studentID, orgID, componentID).
		First(&menuOrg).Error; err != nil {
		return nil, err
	}
	return &menuOrg, nil
}

func (r *StudentMenuOrganizationRepository) UpdateWithTx(
	tx *gorm.DB,
	menuOrg *entity.StudentMenuOrganization,
) error {
	return tx.Save(menuOrg).Error
}

func (r *StudentMenuOrganizationRepository) CreateWithTx(
	tx *gorm.DB,
	menuOrg *entity.StudentMenuOrganization,
) error {
	return tx.Create(menuOrg).Error
}

func (r *StudentMenuOrganizationRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.StudentMenuOrganization{}).Error
	if err != nil {
		log.Error("DeviceMenuRepository.DeleteByComponentID: " + err.Error())
		return errors.New("failed to delete student menu org by component ID")
	}
	return nil
}
