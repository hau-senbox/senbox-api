package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type OrganizationSettingMenuRepository struct {
	DBConn *gorm.DB
}

func (r *OrganizationSettingMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.OrganizationSettingMenu) error {
	return tx.Create(menu).Error
}

func (r *OrganizationSettingMenuRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.OrganizationSettingMenu{}).Error
	if err != nil {
		return errors.New("failed to delete organization setting menu by component ID")
	}
	return nil
}

func (r *OrganizationSettingMenuRepository) GetByOrganizationSettingIDAndComponentID(tx *gorm.DB, organizationSettingID string, componentID string) (*entity.OrganizationSettingMenu, error) {
	var menu entity.OrganizationSettingMenu
	err := tx.
		Where("organization_setting_id = ? AND component_id = ?", organizationSettingID, componentID).
		First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &menu, nil
}

func (r *OrganizationSettingMenuRepository) GetByOrganiazationSettingID(organizationSettingID string) ([]*entity.OrganizationSettingMenu, error) {
	var menus []*entity.OrganizationSettingMenu
	err := r.DBConn.
		Where("organization_setting_id = ?", organizationSettingID).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}
