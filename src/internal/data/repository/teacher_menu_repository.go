package repository

import (
	"errors"
	"fmt"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type TeacherMenuRepository struct {
	DBConn *gorm.DB
}

func NewTeacherMenuRepository(dbConn *gorm.DB) *TeacherMenuRepository {
	return &TeacherMenuRepository{DBConn: dbConn}
}

// Create single teacher menu
func (r *TeacherMenuRepository) Create(menu *entity.TeacherMenu) error {
	return r.DBConn.Create(menu).Error
}

// Bulk create
func (r *TeacherMenuRepository) BulkCreate(menus []entity.TeacherMenu) error {
	return r.DBConn.Create(&menus).Error
}

// Delete all teacher menu by teacher ID
func (r *TeacherMenuRepository) DeleteByTeacherID(teacherID string) error {
	return r.DBConn.Where("teacher_id = ?", teacherID).Delete(&entity.TeacherMenu{}).Error
}

// Get all teacher menu by teacher ID
func (r *TeacherMenuRepository) GetByTeacherID(teacherID string) ([]entity.TeacherMenu, error) {
	var result []entity.TeacherMenu
	err := r.DBConn.Where("teacher_id = ?", teacherID).Find(&result).Error
	return result, err
}

// Create with transaction
func (r *TeacherMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.TeacherMenu) error {
	return tx.Create(menu).Error
}

// Delete all records (dangerous operation)
func (r *TeacherMenuRepository) DeleteAll() error {
	return r.DBConn.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.TeacherMenu{}).Error
}

// Update is_show field by teacher_id and component_id
func (r *TeacherMenuRepository) UpdateIsShowByTeacherAndComponentID(teacherID, componentID string, isShow bool) error {
	return r.DBConn.Model(&entity.TeacherMenu{}).
		Where("teacher_id = ? AND component_id = ?", teacherID, componentID).
		Update("is_show", isShow).Error
}

// Delete all with transaction
func (r *TeacherMenuRepository) DeleteAllTx(tx *gorm.DB) error {
	if err := tx.Exec("DELETE FROM teacher_menu").Error; err != nil {
		log.Error("TeacherMenuRepository.DeleteAllTx: " + err.Error())
		return fmt.Errorf("Delete all teacher_menu fail: %w", err)
	}
	return nil
}

// Get by teacher ID and component ID
func (r *TeacherMenuRepository) GetByTeacherIDAndComponentID(tx *gorm.DB, teacherID, componentID uuid.UUID) (*entity.TeacherMenu, error) {
	var menu entity.TeacherMenu
	err := tx.
		Where("teacher_id = ? AND component_id = ?", teacherID, componentID).
		First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &menu, nil
}

// Update with transaction
func (r *TeacherMenuRepository) UpdateWithTx(tx *gorm.DB, menu *entity.TeacherMenu) error {
	return tx.Model(&entity.TeacherMenu{}).
		Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"order":   menu.Order,
			"visible": menu.Visible,
			"is_show": menu.IsShow,
		}).Error
}

// Delete by component ID
func (r *TeacherMenuRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.TeacherMenu{}).Error
	if err != nil {
		log.Error("TeacherMenuRepository.DeleteByComponentID: " + err.Error())
		return errors.New("failed to delete teacher menu by component ID")
	}
	return nil
}
