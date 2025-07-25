package repository

import (
	"fmt"
	"sen-global-api/internal/domain/entity"

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
