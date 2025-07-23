package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChildMenuRepository struct {
	DBConn *gorm.DB
}

func NewChildMenuRepository(dbConn *gorm.DB) *ChildMenuRepository {
	return &ChildMenuRepository{DBConn: dbConn}
}

func (r *ChildMenuRepository) Create(menu *entity.ChildMenu) error {
	return r.DBConn.Create(menu).Error
}

func (r *ChildMenuRepository) BulkCreate(menus []entity.ChildMenu) error {
	return r.DBConn.Create(&menus).Error
}

func (r *ChildMenuRepository) DeleteByChildID(childID string) error {
	return r.DBConn.Where("child_id = ?", childID).Delete(&entity.ChildMenu{}).Error
}

func (r *ChildMenuRepository) GetByChildID(childID string) ([]entity.ChildMenu, error) {
	var result []entity.ChildMenu
	err := r.DBConn.Where("child_id = ?", childID).Find(&result).Error
	return result, err
}

func (r *ChildMenuRepository) DeleteAll() error {
	return r.DBConn.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.ChildMenu{}).Error
}

func (r *ChildMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.ChildMenu) error {
	return tx.Create(menu).Error
}

func (r *ChildMenuRepository) Update(menu *entity.ChildMenu) error {
	return r.DBConn.Model(&entity.ChildMenu{}).
		Where("id = ?", menu.ID).
		Updates(menu).Error
}

func (r *ChildMenuRepository) UpdateIsShowByChildAndComponentID(childID, componentID string, isShow bool) error {
	return r.DBConn.Model(&entity.ChildMenu{}).
		Where("child_id = ? AND component_id = ?", childID, componentID).
		Update("is_show", isShow).Error
}

func (r *ChildMenuRepository) GetByChildIDAndComponentID(tx *gorm.DB, childID, componentID uuid.UUID) (*entity.ChildMenu, error) {
	var menu entity.ChildMenu
	err := tx.
		Where("child_id = ? AND component_id = ?", childID, componentID).
		First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &menu, nil
}

func (r *ChildMenuRepository) UpdateWithTx(tx *gorm.DB, menu *entity.ChildMenu) error {
	return tx.Model(&entity.ChildMenu{}).
		Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"order":   menu.Order,
			"visible": menu.Visible,
			"is_show": menu.IsShow,
		}).Error
}

func (r *ChildMenuRepository) GetByChildIDActive(childID string) ([]entity.ChildMenu, error) {
	var result []entity.ChildMenu
	err := r.DBConn.Where("child_id = ? AND is_show = ?", childID, true).Find(&result).Error
	return result, err
}
