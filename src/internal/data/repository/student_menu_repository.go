package repository

import (
	"errors"
	"fmt"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StudentMenuRepository struct {
	DBConn *gorm.DB
}

func NewStudentMenuRepository(dbConn *gorm.DB) *StudentMenuRepository {
	return &StudentMenuRepository{DBConn: dbConn}
}

func (r *StudentMenuRepository) Create(menu *entity.StudentMenu) error {
	return r.DBConn.Create(menu).Error
}

func (r *StudentMenuRepository) BulkCreate(menus []entity.StudentMenu) error {
	return r.DBConn.Create(&menus).Error
}

func (r *StudentMenuRepository) DeleteByStudentID(studentID string) error {
	return r.DBConn.Where("student_id = ?", studentID).Delete(&entity.StudentMenu{}).Error
}

func (r *StudentMenuRepository) GetByStudentID(studentID string) ([]entity.StudentMenu, error) {
	var result []entity.StudentMenu
	err := r.DBConn.Where("student_id = ?", studentID).Find(&result).Error
	return result, err
}

func (r *StudentMenuRepository) GetByStudentIDActive(studentID string) ([]entity.StudentMenu, error) {
	var result []entity.StudentMenu
	err := r.DBConn.Where("student_id = ? AND is_show = ?", studentID, true).Find(&result).Error
	return result, err
}

func (r *StudentMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.StudentMenu) error {
	return tx.Create(menu).Error
}

func (r *StudentMenuRepository) DeleteAll() error {
	return r.DBConn.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.StudentMenu{}).Error
}

func (r *StudentMenuRepository) UpdateIsShowByStudentAndComponentID(studentID, componentID string, isShow bool) error {
	return r.DBConn.Model(&entity.StudentMenu{}).
		Where("student_id = ? AND component_id = ?", studentID, componentID).
		Update("is_show", isShow).Error
}

func (r *StudentMenuRepository) DeleteAllTx(tx *gorm.DB) error {
	if err := tx.Exec("DELETE FROM student_menu").Error; err != nil {
		log.Error("StudentMenuRepository.DeleteAllTx: " + err.Error())
		return fmt.Errorf("Delete all student_menu fail: %w", err)
	}
	return nil
}

func (r *StudentMenuRepository) GetByStudentIDAndComponentID(tx *gorm.DB, studentID, componentID uuid.UUID) (*entity.StudentMenu, error) {
	var menu entity.StudentMenu
	err := tx.
		Where("student_id = ? AND component_id = ?", studentID, componentID).
		First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &menu, nil
}

func (r *StudentMenuRepository) UpdateWithTx(tx *gorm.DB, menu *entity.StudentMenu) error {
	return tx.Model(&entity.StudentMenu{}).
		Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"order":   menu.Order,
			"visible": menu.Visible,
			"is_show": menu.IsShow,
		}).Error
}

func (r *StudentMenuRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.StudentMenu{}).Error
	if err != nil {
		log.Error("StudentMenuRepository.DeleteByComponentID: " + err.Error())
		return errors.New("failed to delete student menu by component ID")
	}
	return nil
}
