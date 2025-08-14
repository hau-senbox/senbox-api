package repository

import (
	"sen-global-api/internal/domain/entity"

	"gorm.io/gorm"
)

type OrganizationSettingRepository struct {
	DBConn *gorm.DB
}

func NewOrganizationSettingRepository(db *gorm.DB) *OrganizationSettingRepository {
	return &OrganizationSettingRepository{DBConn: db}
}

// Create inserts a new OrganizationSetting record
func (r *OrganizationSettingRepository) Create(setting *entity.OrganizationSetting) error {
	return r.DBConn.Create(setting).Error
}

// GetByID retrieves an OrganizationSetting by ID
func (r *OrganizationSettingRepository) GetByID(id uint) (*entity.OrganizationSetting, error) {
	var setting entity.OrganizationSetting
	if err := r.DBConn.First(&setting, id).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// Update modifies an existing OrganizationSetting
func (r *OrganizationSettingRepository) Update(setting *entity.OrganizationSetting) error {
	return r.DBConn.Save(setting).Error
}

// Delete removes an OrganizationSetting by ID
func (r *OrganizationSettingRepository) Delete(id uint) error {
	return r.DBConn.Delete(&entity.OrganizationSetting{}, id).Error
}

// List retrieves all OrganizationSettings
func (r *OrganizationSettingRepository) List() ([]entity.OrganizationSetting, error) {
	var settings []entity.OrganizationSetting
	if err := r.DBConn.Find(&settings).Error; err != nil {
		return nil, err
	}
	return settings, nil
}

// GetByOrgID lấy thông tin OrganizationSetting theo OrganizationID
func (r *OrganizationSettingRepository) GetByOrgID(orgID string) (*entity.OrganizationSetting, error) {
	var setting entity.OrganizationSetting
	if err := r.DBConn.Where("organization_id = ?", orgID).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// GetByOrgID lấy thông tin OrganizationSetting theo OrganizationID
func (r *OrganizationSettingRepository) GetByDeviceIdAndOrgId(deviceID string, orgID string) (*entity.OrganizationSetting, error) {
	var setting entity.OrganizationSetting
	if err := r.DBConn.Where("organization_id = ? AND device_id = ?", orgID, deviceID).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// GetByOrgID lấy thông tin OrganizationSetting theo deviceID
func (r *OrganizationSettingRepository) GetByDeviceID(deviceID string) (*entity.OrganizationSetting, error) {
	var setting entity.OrganizationSetting
	if err := r.DBConn.Where("device_id = ?", deviceID).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}

// CreateWithTx tạo mới OrganizationSetting trong transaction
func (r *OrganizationSettingRepository) CreateWithTx(tx *gorm.DB, setting *entity.OrganizationSetting) error {
	return tx.Create(setting).Error
}

// UpdateWithTx cập nhật OrganizationSetting trong transaction
func (r *OrganizationSettingRepository) UpdateWithTx(tx *gorm.DB, setting *entity.OrganizationSetting) error {
	return tx.Save(setting).Error
}

// GetByID retrieves an OrganizationSetting by ID
func (r *OrganizationSettingRepository) GetByIDAndIsNewConfig(orgID string) (*entity.OrganizationSetting, error) {
	var setting entity.OrganizationSetting
	if err := r.DBConn.Where("organization_id = ? AND is_news_config = ?", orgID, true).First(&setting).Error; err != nil {
		return nil, err
	}
	return &setting, nil
}
