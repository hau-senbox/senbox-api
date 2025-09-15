package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type SuperAdminEmergencyMenuRepository struct {
	DBConn *gorm.DB
}

func (r *SuperAdminEmergencyMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.SuperAdminEmergencyMenu) error {
	return tx.Create(menu).Error
}

func (r *SuperAdminEmergencyMenuRepository) UpdateWithTx(tx *gorm.DB, menu *entity.SuperAdminEmergencyMenu) error {
	return tx.Model(&entity.SuperAdminEmergencyMenu{}).
		Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"order": menu.Order,
		}).Error
}

func (r *SuperAdminEmergencyMenuRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.SuperAdminEmergencyMenu{}).Error
	if err != nil {
		return errors.New("failed to delete department menu by component ID")
	}
	return nil
}

func (r *SuperAdminEmergencyMenuRepository) GetByComponentID(componentID string) (*entity.SuperAdminEmergencyMenu, error) {
	var menu entity.SuperAdminEmergencyMenu
	err := r.DBConn.Where("component_id = ?", componentID).First(&menu).Error
	if err != nil {
		return nil, err
	}
	return &menu, nil
}

func (r *SuperAdminEmergencyMenuRepository) GetAll() ([]entity.SuperAdminEmergencyMenu, error) {
	var menus []entity.SuperAdminEmergencyMenu
	err := r.DBConn.Find(&menus).Error
	return menus, err
}
