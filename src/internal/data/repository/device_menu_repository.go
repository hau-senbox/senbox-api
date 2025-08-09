package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type DeviceMenuRepository struct {
	DBConn *gorm.DB
}

func NewDeviceMenuRepository(dbConn *gorm.DB) *DeviceMenuRepository {
	return &DeviceMenuRepository{DBConn: dbConn}
}

func (r *DeviceMenuRepository) Create(menu *entity.SDeviceMenuV2) error {
	return r.DBConn.Create(menu).Error
}

func (r *DeviceMenuRepository) BulkCreate(menus []entity.SDeviceMenuV2) error {
	return r.DBConn.Create(&menus).Error
}

func (r *DeviceMenuRepository) DeleteByDeviceID(deviceID string) error {
	return r.DBConn.Where("device_id = ?", deviceID).Delete(&entity.SDeviceMenuV2{}).Error
}

func (r *DeviceMenuRepository) GetByDeviceID(deviceID string) ([]entity.SDeviceMenuV2, error) {
	var result []entity.SDeviceMenuV2
	err := r.DBConn.Where("device_id = ?", deviceID).Find(&result).Error
	return result, err
}

func (r *DeviceMenuRepository) DeleteAll() error {
	return r.DBConn.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.SDeviceMenuV2{}).Error
}

func (r *DeviceMenuRepository) CreateWithTx(tx *gorm.DB, menu *entity.SDeviceMenuV2) error {
	return tx.Create(menu).Error
}

func (r *DeviceMenuRepository) Update(menu *entity.SDeviceMenuV2) error {
	return r.DBConn.Model(&entity.SDeviceMenuV2{}).
		Where("id = ?", menu.ID).
		Updates(menu).Error
}

func (r *DeviceMenuRepository) UpdateIsShowByDeviceAndComponentID(deviceID string, componentID uuid.UUID, isShow bool) error {
	return r.DBConn.Model(&entity.SDeviceMenuV2{}).
		Where("device_id = ? AND component_id = ?", deviceID, componentID).
		Update("is_show", isShow).Error
}

func (r *DeviceMenuRepository) GetByDeviceIDAndComponentID(tx *gorm.DB, deviceID string, componentID uuid.UUID) (*entity.SDeviceMenuV2, error) {
	var menu entity.SDeviceMenuV2
	err := tx.
		Where("device_id = ? AND component_id = ?", deviceID, componentID).
		First(&menu).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, gorm.ErrRecordNotFound
		}
		return nil, err
	}
	return &menu, nil
}

func (r *DeviceMenuRepository) UpdateWithTx(tx *gorm.DB, menu *entity.SDeviceMenuV2) error {
	return tx.Model(&entity.SDeviceMenuV2{}).
		Where("id = ?", menu.ID).
		Updates(map[string]interface{}{
			"order":   menu.Order,
			"visible": menu.Visible,
			"is_show": menu.IsShow,
		}).Error
}

func (r *DeviceMenuRepository) GetByDeviceIDActive(deviceID string) ([]entity.SDeviceMenuV2, error) {
	var result []entity.SDeviceMenuV2
	err := r.DBConn.Where("device_id = ? AND is_show = ?", deviceID, true).Find(&result).Error
	return result, err
}

func (r *DeviceMenuRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.SDeviceMenuV2{}).Error
	if err != nil {
		log.Error("DeviceMenuRepository.DeleteByComponentID: " + err.Error())
		return errors.New("failed to delete device menu by component ID")
	}
	return nil
}
