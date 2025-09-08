package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type DepartmentMenuRepository struct {
	DBConn *gorm.DB
}

func (r *DepartmentMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.DepartmentMenu) error {
	return tx.Create(menu).Error
}

func (r *DepartmentMenuRepository) UpdateWithTx(tx *gorm.DB, menu *entity.DepartmentMenu) error {
	return tx.Model(&entity.DepartmentMenu{}).
		Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"order": menu.Order,
		}).Error
}

func (r *DepartmentMenuRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.DepartmentMenu{}).Error
	if err != nil {
		return errors.New("failed to delete department menu by component ID")
	}
	return nil
}

func (r *DepartmentMenuRepository) GetByDepartmentIDAndComponentID(tx *gorm.DB, departmentID string, componentID string) (*entity.DepartmentMenu, error) {
	var menu entity.DepartmentMenu
	err := tx.
		Where("department_id = ? AND component_id = ?", departmentID, componentID).
		First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &menu, nil
}

func (r *DepartmentMenuRepository) GetByDepartmentID(departmentID string) ([]*entity.DepartmentMenu, error) {
	var menus []*entity.DepartmentMenu
	err := r.DBConn.
		Where("department_id = ?", departmentID).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}
