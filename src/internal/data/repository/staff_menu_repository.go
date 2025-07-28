package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type StaffMenuRepository struct {
	DBConn *gorm.DB
}

func NewStaffMenuRepository(dbConn *gorm.DB) *StaffMenuRepository {
	return &StaffMenuRepository{DBConn: dbConn}
}

// Create single staff menu
func (r *StaffMenuRepository) Create(menu *entity.StaffMenu) error {
	return r.DBConn.Create(menu).Error
}

// Bulk create
func (r *StaffMenuRepository) BulkCreate(menus []entity.StaffMenu) error {
	return r.DBConn.Create(&menus).Error
}

// Delete all staff menu by staff ID
func (r *StaffMenuRepository) DeleteByStaffID(staffID string) error {
	return r.DBConn.Where("staff_id = ?", staffID).Delete(&entity.StaffMenu{}).Error
}

// Get all staff menu by staff ID
func (r *StaffMenuRepository) GetByStaffID(staffID string) ([]entity.StaffMenu, error) {
	var result []entity.StaffMenu
	err := r.DBConn.Where("staff_id = ?", staffID).Find(&result).Error
	return result, err
}

func (r *StaffMenuRepository) GetByStaffIDActive(staffID string) ([]entity.StaffMenu, error) {
	var result []entity.StaffMenu
	err := r.DBConn.Where("staff_id = ? AND is_show = ?", staffID, true).Find(&result).Error
	return result, err
}

// Create with transaction
func (r *StaffMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.StaffMenu) error {
	return tx.Create(menu).Error
}

// Delete all records (dangerous operation)
func (r *StaffMenuRepository) DeleteAll() error {
	return r.DBConn.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.StaffMenu{}).Error
}

// Update is_show field by staff_id and component_id
func (r *StaffMenuRepository) UpdateIsShowByStaffAndComponentID(staffID, componentID string, isShow bool) error {
	return r.DBConn.Model(&entity.StaffMenu{}).
		Where("staff_id = ? AND component_id = ?", staffID, componentID).
		Update("is_show", isShow).Error
}

// Get by staff ID and component ID
func (r *StaffMenuRepository) GetByStaffIDAndComponentID(tx *gorm.DB, staffID, componentID uuid.UUID) (*entity.StaffMenu, error) {
	var menu entity.StaffMenu
	err := tx.
		Where("staff_id = ? AND component_id = ?", staffID, componentID).
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
func (r *StaffMenuRepository) UpdateWithTx(tx *gorm.DB, menu *entity.StaffMenu) error {
	return tx.Model(&entity.StaffMenu{}).
		Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"order":   menu.Order,
			"visible": menu.Visible,
			"is_show": menu.IsShow,
		}).Error
}

// Delete by component ID
func (r *StaffMenuRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.StaffMenu{}).Error
	if err != nil {
		log.Error("StaffMenuRepository.DeleteByComponentID: " + err.Error())
		return errors.New("failed to delete staff menu by component ID")
	}
	return nil
}
