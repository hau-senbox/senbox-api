package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type OrganizationRepository struct {
	DBConn *gorm.DB
}

func NewOrganizationRepository(dbConn *gorm.DB) *OrganizationRepository {
	return &OrganizationRepository{DBConn: dbConn}
}

func (receiver *OrganizationRepository) GetAll() ([]*entity.SOrganization, error) {
	var organizations []*entity.SOrganization
	err := receiver.DBConn.Find(&organizations).Error
	if err != nil {
		log.Error("OrganizationRepository.GetAll: " + err.Error())
		return nil, errors.New("failed to get organization")
	}
	return organizations, nil
}

func (receiver *OrganizationRepository) GetByID(id uint) (*entity.SOrganization, error) {
	var organization entity.SOrganization
	err := receiver.DBConn.Where("id = ?", id).First(&organization).Error
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

	// update user organization
	err = receiver.DBConn.Model(&entity.SUserEntity{}).
		Where("id = ?", req.UserId).
		Updates(map[string]interface{}{
			"organization_id": req.OrganizationId,
		}).Error

	if err != nil {
		log.Error("OrganizationRepository.UserJoinOrganization: " + err.Error())
		return errors.New("failed to update user organization")
	}

	return nil
}

func (receiver *OrganizationRepository) GetAllUserByOrganization(organizationID uint) ([]*entity.SUserEntity, error) {
	var users []*entity.SUserEntity
	err := receiver.DBConn.Raw("SELECT * FROM s_user_entity sue WHERE sue.id IN (SELECT suo.user_id FROM s_user_organizations suo WHERE organization_id = ?)", organizationID).
		Scan(&users).Error

	if err != nil {
		log.Error("OrganizationRepository.GetAllUserByOrganization: " + err.Error())
		return nil, errors.New("failed to get all users in this organization")
	}

	return users, nil
}
