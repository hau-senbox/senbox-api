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
