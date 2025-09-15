package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type OrganizationEmergencyMenuRepository struct {
	DBConn *gorm.DB
}

func (r *OrganizationEmergencyMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.OrganizationEmergencyMenu) error {
	return tx.Create(menu).Error
}

func (r *OrganizationEmergencyMenuRepository) UpdateWithTx(tx *gorm.DB, menu *entity.OrganizationEmergencyMenu) error {
	return tx.Model(&entity.OrganizationEmergencyMenu{}).
		Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"order": menu.Order,
		}).Error
}

func (r *OrganizationEmergencyMenuRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.OrganizationEmergencyMenu{}).Error
	if err != nil {
		return errors.New("failed to delete department menu by component ID")
	}
	return nil
}

func (r *OrganizationEmergencyMenuRepository) GetByComponentID(componentID string) (*entity.OrganizationEmergencyMenu, error) {
	var menu entity.OrganizationEmergencyMenu
	err := r.DBConn.Where("component_id = ?", componentID).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *OrganizationEmergencyMenuRepository) GetByOrganizationIDAndComponentID(organizationID string, componentID string) (*entity.OrganizationEmergencyMenu, error) {
	var menu entity.OrganizationEmergencyMenu
	err := r.DBConn.Where("organization_id = ? AND component_id = ?", organizationID, componentID).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *OrganizationEmergencyMenuRepository) GetByOrganizationID(organizationID string) ([]entity.OrganizationEmergencyMenu, error) {
	var menus []entity.OrganizationEmergencyMenu
	err := r.DBConn.Where("organization_id = ?", organizationID).Find(&menus).Error
	return menus, err
}
