package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FunctionClaimPermissionRepository struct {
	DBConn *gorm.DB
}

func NewFunctionClaimPermissionRepository(dbConn *gorm.DB) *FunctionClaimPermissionRepository {
	return &FunctionClaimPermissionRepository{DBConn: dbConn}
}

func (receiver *FunctionClaimPermissionRepository) GetAllByFunctionClaim(functionClaimId int64) ([]entity.SFunctionClaimPermission, error) {
	var permissions []entity.SFunctionClaimPermission
	err := receiver.DBConn.Model(entity.SFunctionClaimPermission{}).Where("function_claim_id = ?", functionClaimId).Find(&permissions).Error
	if err != nil {
		log.Error("FunctionClaimPermissionRepository.GetAll: " + err.Error())
		return nil, errors.New("failed to get all permissions")
	}

	return permissions, err
}

func (receiver *FunctionClaimPermissionRepository) GetByID(req request.GetFunctionClaimPermissionByIdRequest) (*entity.SFunctionClaimPermission, error) {
	var permission entity.SFunctionClaimPermission
	err := receiver.DBConn.Where("id = ?", req.ID).First(&permission).Error
	if err != nil {
		log.Error("FunctionClaimPermissionRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get permission")
	}
	return &permission, nil
}

func (receiver *FunctionClaimPermissionRepository) GetByName(req request.GetFunctionClaimPermissionByNameRequest) (*entity.SFunctionClaimPermission, error) {
	var permission entity.SFunctionClaimPermission
	err := receiver.DBConn.Where("permission_name = ?", req.PermissionName).First(&permission).Error
	if err != nil {
		log.Error("FunctionClaimPermissionRepository.GetByName: " + err.Error())
		return nil, errors.New("failed to get permission")
	}
	return &permission, nil
}

func (receiver *FunctionClaimPermissionRepository) CreateFunctionClaimPermission(req request.CreateFunctionClaimPermissionRequest) error {
	permission, _ := receiver.GetByName(request.GetFunctionClaimPermissionByNameRequest{PermissionName: req.PermissionName})

	if permission != nil {
		log.Error("FunctionClaimPermissionRepository.CreateFunctionClaimPermission: " + permission.PermissionName)
		return errors.New("permission already existed")
	}

	var functionClaimCount int64
	receiver.DBConn.Model(&entity.SFunctionClaim{}).Where("id = ?", req.FunctionClaimId).Count(&functionClaimCount)

	if functionClaimCount == 0 {
		log.Error("FunctionClaimPermissionRepository.CreateFunctionClaimPermission: " + "function claim not found")
		return errors.New("function claim not found")
	}

	permissionReq := entity.SFunctionClaimPermission{
		PermissionName:  req.PermissionName,
		FunctionClaimId: req.FunctionClaimId,
	}
	permissionResult := receiver.DBConn.Create(&permissionReq)

	if permissionResult.Error != nil {
		log.Error("FunctionClaimPermissionRepository.CreateFunctionClaimPermission: " + permissionResult.Error.Error())
		return errors.New("failed to create permission")
	}

	return nil
}

func (receiver *FunctionClaimPermissionRepository) UpdateFunctionClaimPermission(req request.UpdateFunctionClaimPermissionRequest) error {
	updateResult := receiver.DBConn.Model(&entity.SFunctionClaimPermission{}).Where("id = ?", req.PermissionId).
		Updates(map[string]interface{}{
			"permission_name":   req.PermissionName,
			"function_claim_id": req.FunctionClaimId,
		})

	if updateResult.Error != nil {
		log.Error("FunctionClaimPermissionRepository.UpdateFunctionClaimPermission: " + updateResult.Error.Error())
		return errors.New("failed to update permission")
	}

	return nil
}

func (receiver *FunctionClaimPermissionRepository) DeleteFunctionClaimPermission(req request.DeleteFunctionClaimPermissionRequest) error {
	result := receiver.DBConn.Delete(&entity.SFunctionClaimPermission{}, req.ID)
	if result.Error != nil {
		log.Error("FunctionClaimPermissionRepository.DeleteFunctionClaimPermission: " + result.Error.Error())
		return errors.New("failed to delete permission")
	}
	return nil
}
