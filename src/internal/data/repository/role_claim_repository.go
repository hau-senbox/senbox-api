package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type RoleClaimRepository struct {
	DBConn *gorm.DB
}

func NewRoleClaimRepository(dbConn *gorm.DB) *RoleClaimRepository {
	return &RoleClaimRepository{DBConn: dbConn}
}

func (receiver *RoleClaimRepository) GetAll() ([]entity.SRoleClaim, error) {
	var roleClaims []entity.SRoleClaim
	err := receiver.DBConn.Table("s_role_claim").Find(&roleClaims).Error
	if err != nil {
		log.Error("RoleClaimRepository.GetAll: " + err.Error())
		return nil, errors.New("failed to get all role claims")
	}

	return roleClaims, err
}

func (receiver *RoleClaimRepository) GetAllByRole(req request.GetAllRoleClaimByRoleRequest) ([]entity.SRoleClaim, error) {
	var roleClaims []entity.SRoleClaim
	err := receiver.DBConn.Table("s_role_claim").Where("role_id = ?", req.RoleId).Find(&roleClaims).Error
	if err != nil {
		log.Error("RoleClaimRepository.GetAllByRole: " + err.Error())
		return nil, errors.New("failed to get all role claims")
	}

	return roleClaims, err
}

func (receiver *RoleClaimRepository) GetByID(req request.GetRoleClaimByIdRequest) (*entity.SRoleClaim, error) {
	var roleClaim entity.SRoleClaim
	err := receiver.DBConn.Where("id = ?", req.ID).First(&roleClaim).Error
	if err != nil {
		log.Error("RoleClaimRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get role claim")
	}
	return &roleClaim, nil
}

func (receiver *RoleClaimRepository) GetByName(req request.GetRoleClaimByNameRequest) (*entity.SRoleClaim, error) {
	var roleClaim entity.SRoleClaim
	err := receiver.DBConn.Where("role_name = ?", req.ClaimName).First(&roleClaim).Error
	if err != nil {
		log.Error("RoleClaimRepository.GetByName: " + err.Error())
		return nil, errors.New("failed to get role claim")
	}
	return &roleClaim, nil
}

func (receiver *RoleClaimRepository) CreateRoleClaim(req request.CreateRoleClaimRequest) error {
	roleClaim, _ := receiver.GetByName(request.GetRoleClaimByNameRequest{ClaimName: req.ClaimName})

	if roleClaim != nil {
		return errors.New("claim already existed")
	}

	var role entity.SRole
	err := receiver.DBConn.Where("id = ?", req.RoleId).First(&role).Error

	if err != nil {
		log.Error("RoleClaimRepository.CreateRoleClaim: " + err.Error())
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("role doesn't exist")
		}
		return errors.New("failed to get role")
	}

	result := receiver.DBConn.Create(&entity.SRoleClaim{
		ClaimName:  req.ClaimName,
		ClaimValue: req.ClaimValue,
		RoleId:     role.ID,
	})

	if result.Error != nil {
		log.Error("RoleClaimRepository.CreateRoleClaim: " + result.Error.Error())
		return errors.New("failed to create role claim")
	}

	return nil
}

func (receiver *RoleClaimRepository) UpdateRoleClaim(req request.UpdateRoleClaimRequest) error {
	updateResult := receiver.DBConn.Model(&entity.SRoleClaim{}).Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"claim_name":  req.ClaimName,
			"claim_value": req.ClaimValue,
		})

	if updateResult.Error != nil {
		log.Error("RoleClaimRepository.UpdateRoleClaim: " + updateResult.Error.Error())
		return errors.New("failed to update role claim")
	}

	return nil
}

func (receiver *RoleClaimRepository) DeleteRoleClaim(req request.DeleteRoleClaimRequest) error {
	deleteResult := receiver.DBConn.Delete(&entity.SRoleClaim{}, req.ID).Error

	if deleteResult != nil {
		log.Error("RoleClaimRepository.DeleteRoleClaim: " + deleteResult.Error())
		return errors.New("failed to delete role claim")
	}
	return nil
}
