package repository

import (
	"fmt"
	"sen-global-api/internal/domain/entity"

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
		return fmt.Errorf("xóa tất cả student_menu thất bại: %w", err)
	}
	return nil
}
