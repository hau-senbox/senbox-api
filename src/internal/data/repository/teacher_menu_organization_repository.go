package repository

import (
	"context"
	"errors"
	"sen-global-api/internal/domain/entity"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TeacherMenuOrganizationRepository struct {
	DBConn *gorm.DB
}

func (r *TeacherMenuOrganizationRepository) GetByID(ctx context.Context, id string) (*entity.TeacherMenuOrganization, error) {
	var menuOrg entity.TeacherMenuOrganization
	if err := r.DBConn.WithContext(ctx).First(&menuOrg, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &menuOrg, nil
}

func (r *TeacherMenuOrganizationRepository) GetAllByTeacherAndOrg(ctx context.Context, teacherID, orgID string) ([]entity.TeacherMenuOrganization, error) {
	var menuOrgs []entity.TeacherMenuOrganization
	err := r.DBConn.WithContext(ctx).
		Where("teacher_id = ? AND organization_id = ?", teacherID, orgID).
		Order("`order` ASC").
		Find(&menuOrgs).Error
	return menuOrgs, err
}

func (r *TeacherMenuOrganizationRepository) GetByTeacherOrgAndComponentID(
	tx *gorm.DB,
	teacherID, orgID, componentID string,
) (*entity.TeacherMenuOrganization, error) {
	var menuOrg entity.TeacherMenuOrganization
	if err := tx.Where("teacher_id = ? AND organization_id = ? AND component_id = ?", teacherID, orgID, componentID).
		First(&menuOrg).Error; err != nil {
		return nil, err
	}
	return &menuOrg, nil
}

func (r *TeacherMenuOrganizationRepository) UpdateWithTx(
	tx *gorm.DB,
	menuOrg *entity.TeacherMenuOrganization,
) error {
	return tx.Save(menuOrg).Error
}

func (r *TeacherMenuOrganizationRepository) CreateWithTx(
	tx *gorm.DB,
	menuOrg *entity.TeacherMenuOrganization,
) error {
	return tx.Create(menuOrg).Error
}

func (r *TeacherMenuOrganizationRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.TeacherMenuOrganization{}).Error
	if err != nil {
		log.Error("DeviceMenuRepository.DeleteByComponentID: " + err.Error())
		return errors.New("failed to delete teacher menu org by component ID")
	}
	return nil
}
