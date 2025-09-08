package repository

import (
	"context"
	"errors"
	"sen-global-api/internal/domain/entity"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DepartmentMenuOrganizationRepository struct {
	DBConn *gorm.DB
}

func (r *DepartmentMenuOrganizationRepository) GetByID(ctx context.Context, id string) (*entity.DepartmentMenuOrganization, error) {
	var menuOrg entity.DepartmentMenuOrganization
	if err := r.DBConn.WithContext(ctx).First(&menuOrg, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &menuOrg, nil
}

func (r *DepartmentMenuOrganizationRepository) GetAllByDepartmentAndOrg(ctx context.Context, departmentID, orgID string) ([]entity.DepartmentMenuOrganization, error) {
	var menuOrgs []entity.DepartmentMenuOrganization
	err := r.DBConn.WithContext(ctx).
		Where("department_id = ? AND organization_id = ?", departmentID, orgID).
		Order("`order` ASC").
		Find(&menuOrgs).Error
	return menuOrgs, err
}

func (r *DepartmentMenuOrganizationRepository) GetByDepartmentOrgAndComponentID(
	tx *gorm.DB,
	departmentID, orgID, componentID string,
) (*entity.DepartmentMenuOrganization, error) {
	var menuOrg entity.DepartmentMenuOrganization
	if err := tx.Where("department_id = ? AND organization_id = ? AND component_id = ?", departmentID, orgID, componentID).
		First(&menuOrg).Error; err != nil {
		return nil, err
	}
	return &menuOrg, nil
}

func (r *DepartmentMenuOrganizationRepository) UpdateWithTx(
	tx *gorm.DB,
	menuOrg *entity.DepartmentMenuOrganization,
) error {
	return tx.Save(menuOrg).Error
}

func (r *DepartmentMenuOrganizationRepository) CreateWithTx(
	tx *gorm.DB,
	menuOrg *entity.DepartmentMenuOrganization,
) error {
	return tx.Create(menuOrg).Error
}

func (r *DepartmentMenuOrganizationRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.DepartmentMenuOrganization{}).Error
	if err != nil {
		log.Error("DepartmentMenuOrganizationRepository.DeleteByComponentID: " + err.Error())
		return errors.New("failed to delete deparmetn menu org by component ID")
	}
	return nil
}
