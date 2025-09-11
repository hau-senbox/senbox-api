package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type ClassroomMenuRepository struct {
	DBConn *gorm.DB
}

func (r *ClassroomMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.ClassroomMenu) error {
	return tx.Create(menu).Error
}

func (r *ClassroomMenuRepository) UpdateWithTx(tx *gorm.DB, menu *entity.ClassroomMenu) error {
	return tx.Model(&entity.ClassroomMenu{}).
		Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"order": menu.Order,
		}).Error
}

func (r *ClassroomMenuRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.ClassroomMenu{}).Error
	if err != nil {
		return errors.New("failed to delete classroom menu by component ID")
	}
	return nil
}

func (r *ClassroomMenuRepository) GetByClassroomIDAndComponentID(tx *gorm.DB, classroomID string, componentID string) (*entity.ClassroomMenu, error) {
	var menu entity.ClassroomMenu
	err := tx.
		Where("classroom_id = ? AND component_id = ?", classroomID, componentID).
		First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &menu, nil
}

func (r *ClassroomMenuRepository) GetByClassroomID(classroomID string) ([]*entity.ClassroomMenu, error) {
	var menus []*entity.ClassroomMenu
	err := r.DBConn.
		Where("classroom_id = ?", classroomID).
		Find(&menus).Error
	if err != nil {
		return nil, err
	}
	return menus, nil
}
