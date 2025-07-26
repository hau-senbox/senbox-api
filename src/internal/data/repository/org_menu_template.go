package repository

import (
	"errors"
	"fmt"
	"sen-global-api/internal/domain/entity"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrganizationMenuTemplateRepository struct {
	DBConn *gorm.DB
}

func NewOrganizationMenuTemplateRepository(db *gorm.DB) *OrganizationMenuTemplateRepository {
	return &OrganizationMenuTemplateRepository{DBConn: db}
}

// Create single record
func (r *OrganizationMenuTemplateRepository) Create(template *entity.OrganizationMenuTemplate) error {
	return r.DBConn.Create(template).Error
}

// Bulk create
func (r *OrganizationMenuTemplateRepository) BulkCreate(templates []entity.OrganizationMenuTemplate) error {
	return r.DBConn.Create(&templates).Error
}

// Get by organization_id
func (r *OrganizationMenuTemplateRepository) GetByOrganizationID(orgID string) ([]entity.OrganizationMenuTemplate, error) {
	var result []entity.OrganizationMenuTemplate
	err := r.DBConn.Where("organization_id = ?", orgID).Find(&result).Error
	return result, err
}

// Delete by organization_id
func (r *OrganizationMenuTemplateRepository) DeleteByOrganizationID(orgID string) error {
	return r.DBConn.Where("organization_id = ?", orgID).Delete(&entity.OrganizationMenuTemplate{}).Error
}

// Delete all records (dangerous)
func (r *OrganizationMenuTemplateRepository) DeleteAll() error {
	return r.DBConn.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&entity.OrganizationMenuTemplate{}).Error
}

// Delete with transaction
func (r *OrganizationMenuTemplateRepository) DeleteAllTx(tx *gorm.DB) error {
	if err := tx.Exec("DELETE FROM organization_menu_template").Error; err != nil {
		log.Error("OrganizationMenuTemplateRepository.DeleteAllTx: " + err.Error())
		return fmt.Errorf("failed to delete all organization_menu_template: %w", err)
	}
	return nil
}

// Create with transaction
func (r *OrganizationMenuTemplateRepository) CreateWithTx(tx *gorm.DB, template *entity.OrganizationMenuTemplate) error {
	return tx.Create(template).Error
}

// GetByOrgIDComponentIDSectionID returns OrganizationMenuTemplate if exists
func (r *OrganizationMenuTemplateRepository) GetByOrgIDComponentIDSectionID(
	tx *gorm.DB,
	orgID string,
	componentID uuid.UUID,
	sectionID uuid.UUID,
) (*entity.OrganizationMenuTemplate, error) {
	var template entity.OrganizationMenuTemplate
	err := tx.Where("organization_id = ? AND component_id = ? AND section_id = ?",
		orgID, componentID.String(), sectionID.String()).First(&template).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // không có record thì trả về nil
		}
		log.Errorf("GetByOrgIDComponentIDSectionID: %v", err)
		return nil, fmt.Errorf("get OrganizationMenuTemplate failed: %w", err)
	}

	return &template, nil
}

// UpdateWithTx updates an existing OrganizationMenuTemplate inside a transaction
func (r *OrganizationMenuTemplateRepository) UpdateWithTx(tx *gorm.DB, template *entity.OrganizationMenuTemplate) error {
	return tx.Save(template).Error
}

// GetBySectionIDAndOrganizationID returns templates matching both section_id and organization_id
func (r *OrganizationMenuTemplateRepository) GetBySectionIDAndOrganizationID(sectionID string, organizationID string) ([]entity.OrganizationMenuTemplate, error) {
	var templates []entity.OrganizationMenuTemplate
	err := r.DBConn.Where("section_id = ? AND organization_id = ?", sectionID, organizationID).Find(&templates).Error
	return templates, err
}

// Delete by component ID
func (r *OrganizationMenuTemplateRepository) DeleteByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&entity.OrganizationMenuTemplate{}).Error
	if err != nil {
		log.Error("OrganizationMenuTemplateRepository.DeleteByComponentID: " + err.Error())
		return errors.New("failed to delete org template menu by component ID")
	}
	return nil
}
