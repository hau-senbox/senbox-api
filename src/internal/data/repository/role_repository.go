package repository

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"
)

type RoleRepository struct {
	DBConn *gorm.DB
}

func NewRoleRepository(dbConn *gorm.DB) *RoleRepository {
	return &RoleRepository{DBConn: dbConn}
}

func (receiver *RoleRepository) GetAll() ([]entity.SRole, error) {
	var roles []entity.SRole
	err := receiver.DBConn.Table("s_role").Find(&roles).Error
	if err != nil {
		log.Error("RoleRepository.GetAll: " + err.Error())
		return nil, errors.New("failed to get all roles")
	}

	return roles, err
}

func (receiver *RoleRepository) GetByID(req request.GetRoleByIdRequest) (*entity.SRole, error) {
	var userRole entity.SRole
	err := receiver.DBConn.Where("id = ?", req.ID).First(&userRole).Error
	if err != nil {
		log.Error("RoleRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get role")
	}
	return &userRole, nil
}

func (receiver *RoleRepository) GetByName(req request.GetRoleByNameRequest) (*entity.SRole, error) {
	var userRole entity.SRole
	err := receiver.DBConn.Where("role_name = ?", req.RoleName).First(&userRole).Error
	if err != nil {
		log.Error("RoleRepository.GetByName: " + err.Error())
		return nil, errors.New("failed to get role or role not found")
	}
	return &userRole, nil
}

func (receiver *RoleRepository) CreateRole(req request.CreateRoleRequest) error {
	userRole, _ := receiver.GetByName(request.GetRoleByNameRequest{RoleName: req.RoleName})

	if userRole != nil {
		return errors.New("role already existed")
	}

	result := receiver.DBConn.Create(&entity.SRole{
		RoleName: req.RoleName,
	})

	if result.Error != nil {
		log.Error("RoleRepository.CreateRole: " + result.Error.Error())
		return errors.New("failed to create role")
	}

	return nil
}

func (receiver *RoleRepository) UpdateRole(req request.UpdateRoleRequest) error {
	updateResult := receiver.DBConn.Model(&entity.SRole{}).Where("id = ?", req.RoleId).
		Updates(map[string]interface{}{
			"role_name": req.RoleName,
		})

	if updateResult.Error != nil {
		log.Error("RoleRepository.UpdateRole: " + updateResult.Error.Error())
		return errors.New("failed to update role")
	}

	return nil
}

func (receiver *RoleRepository) DeleteRole(req request.DeleteRoleRequest) error {
	deleteResult := receiver.DBConn.Delete(&entity.SRole{}, req.ID)

	if deleteResult.Error != nil {
		log.Error("RoleRepository.DeleteRole: " + deleteResult.Error.Error())
		return errors.New("failed to delete role")
	}
	return nil
}
