package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type ParentMenuRepository struct {
	DBConn *gorm.DB
}

func NewParentMenuRepository(dbConn *gorm.DB) *ParentMenuRepository {
	return &ParentMenuRepository{DBConn: dbConn}
}

func (r *ParentMenuRepository) Create(menu *entity.ParentMenu) error {
	return r.DBConn.Create(menu).Error
}

func (r *ParentMenuRepository) BulkCreate(menus []entity.ParentMenu) error {
	return r.DBConn.Create(&menus).Error
}

func (r *ParentMenuRepository) DeleteByParentID(parentID string) error {
	return r.DBConn.Where("parent_id = ?", parentID).Delete(&entity.ParentMenu{}).Error
}

func (r *ParentMenuRepository) GetByParentID(parentID string) ([]entity.ParentMenu, error) {
	var result []entity.ParentMenu
	err := r.DBConn.Where("parent_id = ?", parentID).Find(&result).Error
	return result, err
}

func (r *ParentMenuRepository) DeleteAll() error {
	return r.DBConn.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.ParentMenu{}).Error
}

func (r *ParentMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.ParentMenu) error {
	return tx.Create(menu).Error
}

func (r *ParentMenuRepository) Update(menu *entity.ParentMenu) error {
	return r.DBConn.Model(&entity.ParentMenu{}).
		Where("id = ?", menu.ID).
		Updates(menu).Error
}

func (r *ParentMenuRepository) UpdateIsShowByParentAndComponentID(parentID, componentID string, isShow bool) error {
	return r.DBConn.Model(&entity.ParentMenu{}).
		Where("parent_id = ? AND component_id = ?", parentID, componentID).
		Update("is_show", isShow).Error
}

func (r *ParentMenuRepository) GetByParentIDAndComponentID(tx *gorm.DB, parentID, componentID uuid.UUID) (*entity.ParentMenu, error) {
	var menu entity.ParentMenu
	err := tx.
		Where("parent_id = ? AND component_id = ?", parentID, componentID).
		First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &menu, nil
}

func (r *ParentMenuRepository) UpdateWithTx(tx *gorm.DB, menu *entity.ParentMenu) error {
	return tx.Model(&entity.ParentMenu{}).
		Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"order":   menu.Order,
			"visible": menu.Visible,
			"is_show": menu.IsShow,
		}).Error
}

func (r *ParentMenuRepository) GetByParentIDActive(parentID string) ([]entity.ParentMenu, error) {
	var result []entity.ParentMenu
	err := r.DBConn.Where("parent_id = ?", parentID).Find(&result).Error
	return result, err
}

func (r *ParentMenuRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.ParentMenu{}).Error
	if err != nil {
		log.Error("ParentMenuRepository.DeleteByComponentID: " + err.Error())
		return errors.New("failed to delete parent menu by component ID")
	}
	return nil
}
