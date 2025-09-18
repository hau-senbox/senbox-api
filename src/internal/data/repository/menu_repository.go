package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity/menu"
	"sen-global-api/internal/domain/request"

	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type MenuRepository struct {
	DBConn *gorm.DB
}

func NewMenuRepository(dbConn *gorm.DB) *MenuRepository {
	return &MenuRepository{DBConn: dbConn}
}

func (receiver *MenuRepository) GetSuperAdminMenu() ([]menu.SuperAdminMenu, error) {
	var menus []menu.SuperAdminMenu
	err := receiver.DBConn.Model(&menu.SuperAdminMenu{}).
		Preload("Component").
		Find(&menus).Error
	if err != nil {
		log.Error("MenuRepository.GetSuperAdminMenu: " + err.Error())
		return nil, errors.New("failed to get super admin menu")
	}

	return menus, nil
}

func (receiver *MenuRepository) GetOrgMenu(orgID string) ([]menu.OrgMenu, error) {
	var menus []menu.OrgMenu
	err := receiver.DBConn.Model(&menu.OrgMenu{}).
		Where("organization_id = ?", orgID).
		Preload("Component").
		Find(&menus).Error
	if err != nil {
		log.Error("MenuRepository.GetOrgMenu: " + err.Error())
		return nil, errors.New("failed to get org menu")
	}

	return menus, nil
}

func (receiver *MenuRepository) GetUserMenu(userID string) ([]menu.UserMenu, error) {
	var menus []menu.UserMenu

	err := receiver.DBConn.Model(&menu.UserMenu{}).
		Where("user_id = ?", userID).
		Preload("Component").
		Find(&menus).Error
	if err != nil {
		log.Error("MenuRepository.GetUserMenu: " + err.Error())
		return nil, errors.New("failed to get user menu")
	}

	return menus, nil
}

func (receiver *MenuRepository) GetDeviceMenu(deviceID string) ([]menu.DeviceMenu, error) {
	var menus []menu.DeviceMenu

	err := receiver.DBConn.Model(&menu.DeviceMenu{}).
		Preload("Component").
		Joins("INNER JOIN s_organization o ON o.id = device_menu.organization_id").
		Joins("INNER JOIN s_org_devices od ON od.organization_id = o.id").
		Where("od.device_id = ?", deviceID).
		Find(&menus).Error
	if err != nil {
		log.Error("MenuRepository.GetDeviceMenu: " + err.Error())
		return nil, errors.New("failed to get device menu")
	}

	return menus, nil
}

func (receiver *MenuRepository) GetDeviceMenuByOrg(organizationID string) ([]menu.DeviceMenu, error) {
	var menus []menu.DeviceMenu

	err := receiver.DBConn.Model(&menu.DeviceMenu{}).
		Preload("Component").
		Where("organization_id = ?", organizationID).
		Find(&menus).Error
	if err != nil {
		log.Error("MenuRepository.GetDeviceMenu: " + err.Error())
		return nil, errors.New("failed to get device menu")
	}

	return menus, nil
}

func (receiver *MenuRepository) CreateSuperAdminMenu(req request.CreateSuperAdminMenuRequest, tx *gorm.DB) error {
	var menus []menu.SuperAdminMenu
	for _, component := range req.Components {
		menus = append(menus, menu.SuperAdminMenu{
			Direction:   req.Direction,
			ComponentID: *component.ID,
			Order:       component.Order,
		})
	}

	if tx == nil {
		err := receiver.DBConn.Create(&menus).Error
		if err != nil {
			log.Error("MenuRepository.CreateSuperAdminMenu: " + err.Error())
			return errors.New("failed to create super admin menu")
		}

		return nil
	}

	err := tx.Create(&menus).Error
	if err != nil {
		tx.Rollback()
		log.Error("MenuRepository.CreateSuperAdminMenu: " + err.Error())
		return errors.New("failed to create super admin menu")
	}

	return nil
}

func (receiver *MenuRepository) DeleteSuperAdminMenu(tx *gorm.DB) error {
	if tx == nil {
		err := receiver.DBConn.Exec("DELETE c, sam FROM component c JOIN super_admin_menu sam ON c.id = sam.component_id WHERE 1").Error
		if err != nil {
			log.Error("MenuRepository.DeleteSuperAdminMenu: " + err.Error())
			return errors.New("failed to delete super admin menu")
		}

		return nil
	}

	err := tx.Exec("DELETE c, sam FROM component c JOIN super_admin_menu sam ON c.id = sam.component_id WHERE 1").Error
	if err != nil {
		tx.Rollback()
		log.Error("MenuRepository.DeleteSuperAdminMenu: " + err.Error())
		return errors.New("failed to delete super admin menu")
	}

	return nil
}

func (receiver *MenuRepository) CreateOrgMenu(req request.CreateOrgMenuRequest, tx *gorm.DB) error {
	var menus []menu.OrgMenu
	for _, component := range req.Components {
		menus = append(menus, menu.OrgMenu{
			OrganizationID: uuid.MustParse(req.OrganizationID),
			Direction:      req.Direction,
			ComponentID:    *component.ID,
			Order:          component.Order,
		})
	}

	if tx == nil {
		err := receiver.DBConn.Create(&menus).Error
		if err != nil {
			log.Error("MenuRepository.CreateOrgMenu: " + err.Error())
			return errors.New("failed to create org menu")
		}

		return nil
	}

	err := tx.Create(&menus).Error
	if err != nil {
		tx.Rollback()
		log.Error("MenuRepository.CreateOrgMenu: " + err.Error())
		return errors.New("failed to create org menu")
	}

	return nil
}

func (receiver *MenuRepository) DeleteOrgMenu(organizationID string, tx *gorm.DB) error {
	if tx == nil {
		err := receiver.DBConn.Exec("DELETE c, om FROM component c JOIN org_menu om ON c.id = om.component_id WHERE om.organization_id = ?", organizationID).Error
		if err != nil {
			log.Error("MenuRepository.DeleteOrgMenu: " + err.Error())
			return errors.New("failed to delete org menu")
		}

		return nil
	}

	err := tx.Exec("DELETE c, om FROM component c JOIN org_menu om ON c.id = om.component_id WHERE om.organization_id = ?", organizationID).Error
	if err != nil {
		tx.Rollback()
		log.Error("MenuRepository.DeleteOrgMenu: " + err.Error())
		return errors.New("failed to delete org menu")
	}

	return nil
}

func (receiver *MenuRepository) CreateUserMenu(req request.CreateUserMenuRequest, tx *gorm.DB) error {
	var menus []menu.UserMenu
	for _, component := range req.Components {
		menus = append(menus, menu.UserMenu{
			UserID:      uuid.MustParse(req.UserID),
			ComponentID: *component.ID,
			Order:       component.Order,
		})
	}

	if tx == nil {
		err := receiver.DBConn.Create(&menus).Error
		if err != nil {
			log.Error("MenuRepository.CreateUserMenu: " + err.Error())
			return errors.New("failed to create user menu")
		}

		return nil
	}

	err := tx.Create(&menus).Error
	if err != nil {
		tx.Rollback()
		log.Error("MenuRepository.CreateUserMenu: " + err.Error())
		return errors.New("failed to create user menu")
	}

	return nil
}

func (receiver *MenuRepository) CreateDeviceMenu(req request.CreateDeviceMenuRequest, tx *gorm.DB) error {
	var menus []menu.DeviceMenu
	for _, component := range req.Components {
		menus = append(menus, menu.DeviceMenu{
			OrganizationID: uuid.MustParse(req.OrganizationID),
			ComponentID:    *component.ID,
			Order:          component.Order,
		})
	}

	if tx == nil {
		err := receiver.DBConn.Create(&menus).Error
		if err != nil {
			log.Error("MenuRepository.CreateDeviceMenu: " + err.Error())
			return errors.New("failed to create device menu")
		}

		return nil
	}

	err := tx.Create(&menus).Error
	if err != nil {
		tx.Rollback()
		log.Error("MenuRepository.CreateDeviceMenu: " + err.Error())
		return errors.New("failed to create device menu")
	}

	return nil
}

func (receiver *MenuRepository) DeleteUserMenu(userID string, tx *gorm.DB) error {
	if tx == nil {
		err := receiver.DBConn.Exec("DELETE c, um FROM component c JOIN user_menu um ON c.id = um.component_id WHERE um.user_id = ?", userID).Error
		if err != nil {
			log.Error("MenuRepository.DeleteUserMenu: " + err.Error())
			return errors.New("failed to delete user menu")
		}

		return nil
	}

	err := tx.Exec("DELETE c, um FROM component c JOIN user_menu um ON c.id = um.component_id WHERE um.user_id = ?", userID).Error
	if err != nil {
		tx.Rollback()
		log.Error("MenuRepository.DeleteUserMenu: " + err.Error())
		return errors.New("failed to delete user menu")
	}

	return nil
}

func (receiver *MenuRepository) DeleteDeviceMenu(organizationID string, tx *gorm.DB) error {
	if tx == nil {
		err := receiver.DBConn.Exec("DELETE c, dm FROM component c JOIN device_menu dm ON c.id = dm.component_id WHERE dm.organization_id = ?", organizationID).Error
		if err != nil {
			log.Error("MenuRepository.DeleteDeviceMenu: " + err.Error())
			return errors.New("failed to delete device menu")
		}

		return nil
	}

	err := tx.Exec("DELETE c, dm FROM component c JOIN device_menu dm ON c.id = dm.component_id WHERE dm.organization_id = ?", organizationID).Error
	if err != nil {
		tx.Rollback()
		log.Error("MenuRepository.DeleteDeviceMenu: " + err.Error())
		return errors.New("failed to delete device menu")
	}

	return nil
}

// Get super admin menu
func (receiver *MenuRepository) GetSuperAdminMenuByComponentID(componentID string) (*menu.SuperAdminMenu, error) {
	var menuItem menu.SuperAdminMenu
	err := receiver.DBConn.
		Where("component_id = ?", componentID).
		First(&menuItem).Error
	if err != nil {
		log.Error("MenuRepository.GetSuperAdminMenuByComponentID: " + err.Error())
		return nil, err
	}

	return &menuItem, nil
}

func (receiver *MenuRepository) UpdateSuperAdminWithTx(tx *gorm.DB, menuItem *menu.SuperAdminMenu) error {
	err := tx.Model(&menu.SuperAdminMenu{}).
		Where("component_id = ?", menuItem.ComponentID).
		Updates(menuItem).Error
	if err != nil {
		log.Error("MenuRepository.UpdateSuperAdminWithTx: " + err.Error())
		return errors.New("failed to update super admin menu by component_id")
	}
	return nil
}

func (receiver *MenuRepository) CreateSuperAdminWithTx(tx *gorm.DB, menuItem *menu.SuperAdminMenu) error {
	err := tx.Create(menuItem).Error
	if err != nil {
		log.Error("MenuRepository.CreateSuperAdminWithTx: " + err.Error())
		return errors.New("failed to create super admin menu")
	}
	return nil
}

func (r *MenuRepository) DeleteSuperAdminMenuByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&menu.SuperAdminMenu{}).Error
	if err != nil {
		log.Error("MenuRepository.DeleteSuperAdminMenuByComponentID: " + err.Error())
		return errors.New("failed to super admin menu by component ID")
	}
	return nil
}

// Organization admin menu
func (receiver *MenuRepository) GetOrganizationAdminMenuByComponentID(componentID string, organizationID string) (*menu.OrgMenu, error) {
	var menuItem menu.OrgMenu
	err := receiver.DBConn.
		Where("component_id = ? AND organization_id = ?", componentID, organizationID).
		First(&menuItem).Error
	if err != nil {
		log.Error("MenuRepository.GetOrganizationAdminMenuByComponentID: " + err.Error())
		return nil, err
	}

	return &menuItem, nil
}

func (receiver *MenuRepository) UpdateOrganizationAdminWithTx(tx *gorm.DB, menuItem *menu.OrgMenu) error {
	err := tx.Model(&menu.OrgMenu{}).
		Where("organization_id = ? AND component_id = ?", menuItem.OrganizationID, menuItem.ComponentID).
		Updates(menuItem).Error
	if err != nil {
		log.Error("MenuRepository.UpdateOrganizationAdminWithTx: " + err.Error())
		return errors.New("failed to update organization admin menu by component_id")
	}
	return nil
}

func (receiver *MenuRepository) CreateOrganizationAdminWithTx(tx *gorm.DB, menuItem *menu.OrgMenu) error {
	err := tx.Create(menuItem).Error
	if err != nil {
		log.Error("MenuRepository.CreateOrganizationAdminWithTx: " + err.Error())
		return errors.New("failed to create organization admin menu")
	}
	return nil
}

func (r *MenuRepository) DeleteOrganizationAdminMenuByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&menu.OrgMenu{}).Error
	if err != nil {
		log.Error("MenuRepository.DeleteOrganizationAdminMenuByComponentID: " + err.Error())
		return errors.New("failed to organization admin menu by component ID")
	}
	return nil
}

func (r *MenuRepository) CreateDeviceMenuOrganization(tx *gorm.DB, deviceMenu *menu.DeviceMenu) error {
	err := tx.Create(deviceMenu).Error
	if err != nil {
		log.Error("MenuRepository.CreateDeviceMenuOrganization: " + err.Error())
		return errors.New("failed to create device menu organization")
	}
	return nil
}

func (r *MenuRepository) DeleteDeviceMenuOrganizationByComponentID(componentID string) error {
	err := r.DBConn.Where("component_id = ?", componentID).Delete(&menu.DeviceMenu{}).Error
	if err != nil {
		log.Error("MenuRepository.DeleteDeviceMenuOrganization: " + err.Error())
		return errors.New("failed to delete device menu organization")
	}
	return nil
}

func (r *MenuRepository) GetDeviceMenuOrganization(organizationID string, componentID string) (*menu.DeviceMenu, error) {
	var deviceMenu menu.DeviceMenu
	err := r.DBConn.
		Where("organization_id = ? AND component_id = ?", organizationID, componentID).
		First(&deviceMenu).Error
	if err != nil {
		log.Error("MenuRepository.GetDeviceMenuOrganization: " + err.Error())
		return nil, err
	}
	return &deviceMenu, nil
}

func (r *MenuRepository) UpdateDeviceMenuOrganizationWithTx(tx *gorm.DB, deviceMenu *menu.DeviceMenu) error {
	err := tx.Model(&menu.DeviceMenu{}).
		Where("organization_id = ? AND component_id = ?", deviceMenu.OrganizationID, deviceMenu.ComponentID).
		Updates(deviceMenu).Error
	if err != nil {
		log.Error("MenuRepository.UpdateDeviceMenuOrganizationWithTx: " + err.Error())
		return errors.New("failed to update device menu organization")
	}
	return nil
}

func (receiver *MenuRepository) GetSuperAdminMenuByLanguage(language uint) ([]menu.SuperAdminMenu, error) {
	var menus []menu.SuperAdminMenu
	err := receiver.DBConn.Model(&menu.SuperAdminMenu{}).
		Preload("Component", "language = ?", language).
		Find(&menus).Error
	if err != nil {
		log.Error("MenuRepository.GetSuperAdminMenuByLanguage: " + err.Error())
		return nil, errors.New("failed to get super admin menu by language")
	}

	return menus, nil
}

func (receiver *MenuRepository) GetOrgMenuByLanguage(orgID string, language uint) ([]menu.OrgMenu, error) {
	var menus []menu.OrgMenu
	err := receiver.DBConn.Model(&menu.OrgMenu{}).
		Where("organization_id = ?", orgID).
		Preload("Component", "language = ?", language).
		Find(&menus).Error
	if err != nil {
		log.Error("MenuRepository.GetOrgMenu: " + err.Error())
		return nil, errors.New("failed to get org menu")
	}

	return menus, nil
}

func (receiver *MenuRepository) GetUserMenuByLanguage(userID string, language uint) ([]menu.UserMenu, error) {
	var menus []menu.UserMenu

	err := receiver.DBConn.Model(&menu.UserMenu{}).
		Where("user_id = ?", userID).
		Preload("Component", "language = ?", language).
		Find(&menus).Error
	if err != nil {
		log.Error("MenuRepository.GetUserMenu: " + err.Error())
		return nil, errors.New("failed to get user menu")
	}

	return menus, nil
}
