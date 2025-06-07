package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/tiendc/gofn"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
	"sen-global-api/internal/domain/value"
	"time"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	DBConn *gorm.DB
}

func NewOrganizationRepository(dbConn *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{DBConn: dbConn}
}

func (receiver *OrganizationRepository) GetAll(user *entity.SUserEntity) ([]*entity.SOrganization, error) {
	var organizations []*entity.SOrganization
	roles := gofn.MapSliceToMap(user.Roles, func(role entity.SRole) (int64, string) {
		return role.ID, role.Role.String()
	})

	if gofn.Contain(gofn.MapValues(roles), "SuperAdmin") {
		err := receiver.DBConn.Preload("UserOrgs").Find(&organizations).Error
		if err != nil {
			log.Error("OrganizationRepository.GetAll: " + err.Error())
			return nil, errors.New("failed to get organization")
		}
	} else {
		var userOrgs []*entity.SUserOrg
		err := receiver.DBConn.Model(&entity.SUserOrg{}).Where("user_id = ?", user.ID).Find(&userOrgs).Error
		if err != nil {
			log.Error("OrganizationRepository.GetAll: " + err.Error())
			return nil, errors.New("failed to get user org")
		}
		orgs := gofn.MapSliceToMap(userOrgs, func(org *entity.SUserOrg) (int64, string) {
			return org.OrganizationId, org.Organization.OrganizationName
		})

		err = receiver.DBConn.Preload("UserOrgs").Where("id in (?)", gofn.MapKeys(orgs)).Find(&organizations).Error
		if err != nil {
			log.Error("OrganizationRepository.GetAll: " + err.Error())
			return nil, errors.New("failed to get organization")
		}
	}

	return organizations, nil
}

func (receiver *OrganizationRepository) GetByID(id int64) (*entity.SOrganization, error) {
	var organization entity.SOrganization
	err := receiver.DBConn.Preload("UserOrgs").Where("id = ?", id).First(&organization).Error
	if err != nil {
		log.Error("OrganizationRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get organization")
	}

	return &organization, nil
}

func (receiver *OrganizationRepository) GetByName(name string) (*entity.SOrganization, error) {
	var organization entity.SOrganization
	err := receiver.DBConn.Preload("UserOrgs").Where("organization_name = ?", name).First(&organization).Error
	if err != nil {
		log.Error("OrganizationRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get organization")
	}

	return &organization, nil
}

func (receiver *OrganizationRepository) CreateOrganization(req request.CreateOrganizationRequest) error {
	result := receiver.DBConn.Create(&entity.SOrganization{
		OrganizationName: req.OrganizationName,
		Password:         req.Password,
		Address:          req.Address,
		Description:      req.Description,
	})

	if result.Error != nil {
		log.Error("OrganizationRepository.CreateOrganization: " + result.Error.Error())
		return errors.New("failed to create organization")
	}

	return nil
}

func (receiver *OrganizationRepository) UpdateOrganization(req request.UpdateOrganizationRequest) error {
	updateResult := receiver.DBConn.Model(&entity.SOrganization{}).Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"organization_name": req.OrganizationName,
			"address":           req.Address,
			"description":       req.Description,
		})

	if updateResult.Error != nil {
		log.Error("OrganizationRepository.UpdateOrganization: " + updateResult.Error.Error())
		return errors.New("failed to update organization")
	}

	return nil
}

func (receiver *OrganizationRepository) UserJoinOrganization(req request.UserJoinOrganizationRequest) error {
	var user entity.SUserEntity
	err := receiver.DBConn.Model(&entity.SUserEntity{}).Where("id = ?", req.UserId).First(&user).Error

	if err != nil {
		log.Error("OrganizationRepository.UserJoinOrganization: " + err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("user doesn't exist")
		}
		return errors.New("failed to get user")
	}

	organization, err := receiver.GetByID(req.OrganizationId)

	if err != nil {
		log.Error("OrganizationRepository.UserJoinOrganization: " + err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("organization doesn't exist")
		}
		return errors.New("failed to get organization")
	}

	var count int64
	err = receiver.DBConn.Table("s_user_organizations").
		Where("user_id = ? AND organization_id = ?", req.UserId, req.OrganizationId).
		Count(&count).Error
	if count > 0 {
		return errors.New("user already joined the organization")
	}

	if err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Error("OrganizationRepository.UserJoinOrganization: " + err.Error())
			return errors.New("failed to fetch user in organization")
		}
	}

	result := receiver.DBConn.Table("s_user_organizations").Create(map[string]interface{}{
		"user_id":         user.ID.String(),
		"organization_id": organization.ID,
	})

	if result.Error != nil {
		log.Error("OrganizationRepository.UserJoinOrganization: " + result.Error.Error())
		return errors.New("failed to assign user for organization")
	}

	return nil
}

func (receiver *OrganizationRepository) GetUserOrgInfo(userId string, organizationId int64) (*entity.SUserOrg, error) {
	var userOrg entity.SUserOrg
	err := receiver.DBConn.Model(&entity.SUserOrg{}).
		Where("user_id = ? AND organization_id = ?", userId, organizationId).
		Find(&userOrg).Error

	if err != nil {
		log.Error("OrganizationRepository.GetUserOrgInfo: " + err.Error())
		return nil, errors.New("failed to get user org info")
	}

	return &userOrg, nil
}

func (receiver *OrganizationRepository) GetAllOrgManagerInfo(organizationId string) (*[]entity.SUserOrg, error) {
	var userOrg []entity.SUserOrg
	err := receiver.DBConn.Model(&entity.SUserOrg{}).
		Where("organization_id = ? AND is_manager = 1", organizationId).
		Find(&userOrg).Error

	if err != nil {
		log.Error("OrganizationRepository.GetAllOrgManagerInfo: " + err.Error())
		return nil, errors.New("failed to get all user org info")
	}

	return &userOrg, nil
}

func (receiver *OrganizationRepository) UpdateUserOrgInfo(req request.UpdateUserOrgInfoRequest) error {
	updateResult := receiver.DBConn.Model(&entity.SUserOrg{}).Where("user_id = ? AND organization_id = ?", req.UserId, req.OrganizationId).
		Updates(map[string]interface{}{
			"user_nick_name": req.UserNickName,
			"is_manager":     req.IsManager,
		})

	if updateResult.Error != nil {
		log.Error("OrganizationRepository.UpdateUserOrgInfo: " + updateResult.Error.Error())
		return errors.New("failed to update user org info")
	}

	return nil
}

func (receiver *OrganizationRepository) GetAllUserByOrganization(organizationID uint) ([]*entity.SUserOrg, error) {
	var userOrg []*entity.SUserOrg
	err := receiver.DBConn.Model(&entity.SUserOrg{}).
		Where("organization_id = ?", organizationID).
		Find(&userOrg).Error

	if err != nil {
		log.Error("OrganizationRepository.GetAllUserByOrganization: " + err.Error())
		return nil, errors.New("failed to get all users in this organization")
	}

	return userOrg, nil
}

func (receiver *OrganizationRepository) GetAllOrgFormApplication() ([]*entity.SOrgFormApplication, error) {
	var forms []*entity.SOrgFormApplication
	err := receiver.DBConn.Model(&entity.SOrgFormApplication{}).Find(&forms).Error

	if err != nil {
		log.Error("OrganizationRepository.GetAllOrgFormApplication: " + err.Error())
		return nil, errors.New("failed to get all application form")
	}

	return forms, nil
}

func (receiver *OrganizationRepository) GetOrgFormApplicationByID(applicationID int64) (*entity.SOrgFormApplication, error) {
	var form entity.SOrgFormApplication
	err := receiver.DBConn.Model(&entity.SOrgFormApplication{}).
		Where("id = ?", applicationID).
		Preload("User").
		First(&form).Error

	if err != nil {
		log.Error("OrganizationRepository.GetAllOrgFormApplication: " + err.Error())
		return nil, errors.New("failed to get application form")
	}

	return &form, nil
}

func (receiver *OrganizationRepository) ApproveOrgFormApplication(applicationID int64) error {
	form, err := receiver.GetOrgFormApplicationByID(applicationID)
	if err != nil {
		log.Error("OrganizationRepository.ApproveOrgFormApplication: " + err.Error())
		return err
	}

	if !form.ApprovedAt.IsZero() {
		return errors.New("organization has already been approved")
	}

	tx := receiver.DBConn.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Update application status
	err = tx.Model(&entity.SOrgFormApplication{}).
		Where("id = ?", applicationID).
		Updates(map[string]interface{}{
			"status":      value.Approved,
			"approved_at": time.Now(),
		}).Error
	if err != nil {
		tx.Rollback()
		log.Error("Update application status failed: " + err.Error())
		return errors.New("failed to approve organization")
	}

	// Create organization
	organization := entity.SOrganization{
		OrganizationName: form.OrganizationName,
	}
	err = tx.Create(&organization).Error
	if err != nil {
		tx.Rollback()
		log.Error("Create organization failed: " + err.Error())
		return errors.New("failed to create organization")
	}

	// Create user organization mapping (Manager)
	userOrg := entity.SUserOrg{
		UserId:         form.UserId,
		OrganizationId: organization.ID,
		UserNickName:   "Manager",
		IsManager:      true,
	}
	err = tx.Create(&userOrg).Error
	if err != nil {
		tx.Rollback()
		log.Error("Create user organization mapping failed: " + err.Error())
		return errors.New("failed to create user-organization relationship")
	}

	// Add user to organization (if needed)
	err = tx.Table("s_user_organizations").Create(map[string]interface{}{
		"user_id":         form.UserId.String(),
		"organization_id": organization.ID,
	}).Error
	if err != nil {
		tx.Rollback()
		log.Error("insert into s_user_organizations failed: " + err.Error())
		return errors.New("failed to join organization")
	}

	// Assign admin role to user
	userRole := entity.SUserRoles{
		UserId: form.UserId,
		RoleId: 5, // admin
	}
	err = tx.Create(&userRole).Error
	if err != nil {
		tx.Rollback()
		log.Error("Assign admin role failed: " + err.Error())
		return errors.New("failed to assign admin role to user")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		log.Error("Transaction commit failed: " + err.Error())
		return errors.New("failed to commit transaction")
	}

	return nil
}

func (receiver *OrganizationRepository) BlockOrgFormApplication(applicationID int64) error {
	form, err := receiver.GetOrgFormApplicationByID(applicationID)

	if err != nil {
		log.Error("OrganizationRepository.ApproveOrgFromApplication: " + err.Error())
		return errors.New("failed to get all application forms")
	}

	if !form.ApprovedAt.IsZero() {
		return errors.New("organization has not been approved")
	}

	err = receiver.DBConn.Model(&entity.SOrgFormApplication{}).Where("id = ?", applicationID).
		Updates(map[string]interface{}{"status": value.Blocked}).Error

	if err != nil {
		log.Error("OrganizationRepository.BlockOrgFromApplication: " + err.Error())
		return errors.New("failed to block organization")
	}

	return nil
}

func (receiver *OrganizationRepository) CreateOrgFormApplication(req request.CreateOrgFormApplicationRequest) error {
	result := receiver.DBConn.Create(&entity.SOrgFormApplication{
		OrganizationName:   req.OrganizationName,
		ApplicationContent: req.ApplicationContent,
		UserId:             uuid.MustParse(req.UserID),
	})

	if result.Error != nil {
		log.Error("OrganizationRepository.CreateOrgFormApplication: " + result.Error.Error())
		return errors.New("failed to create organization form application")
	}

	return nil
}
