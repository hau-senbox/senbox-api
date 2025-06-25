package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type FunctionClaimRepository struct {
	DBConn *gorm.DB
}

func NewFunctionClaimRepository(dbConn *gorm.DB) *FunctionClaimRepository {
	return &FunctionClaimRepository{DBConn: dbConn}
}

func (receiver *FunctionClaimRepository) GetAll() ([]entity.SFunctionClaim, error) {
	var functionClaims []entity.SFunctionClaim
	err := receiver.DBConn.Model(entity.SFunctionClaim{}).
		Preload("ClaimPermissions").
		Find(&functionClaims).Error
	if err != nil {
		log.Error("FunctionClaimRepository.GetAll: " + err.Error())
		return nil, errors.New("failed to get all function claims")
	}

	return functionClaims, err
}

func (receiver *FunctionClaimRepository) GetByID(req request.GetFunctionClaimByIDRequest) (*entity.SFunctionClaim, error) {
	var functionClaim entity.SFunctionClaim
	err := receiver.DBConn.Where("id = ?", req.ID).
		Preload("ClaimPermissions").
		First(&functionClaim).Error
	if err != nil {
		log.Error("FunctionClaimRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get function claim")
	}
	return &functionClaim, nil
}

func (receiver *FunctionClaimRepository) GetByName(req request.GetFunctionClaimByNameRequest) (*entity.SFunctionClaim, error) {
	var functionClaim entity.SFunctionClaim
	err := receiver.DBConn.Where("function_name = ?", req.FunctionName).
		Preload("ClaimPermissions").
		First(&functionClaim).Error
	if err != nil {
		log.Error("FunctionClaimRepository.GetByName: " + err.Error())
		return nil, errors.New("failed to get function claim")
	}
	return &functionClaim, nil
}

func (receiver *FunctionClaimRepository) CreateFunctionClaim(req request.CreateFunctionClaimRequest) error {
	functionClaim, _ := receiver.GetByName(request.GetFunctionClaimByNameRequest{FunctionName: req.FunctionName})

	if functionClaim != nil {
		return errors.New("function claim already existed")
	}

	result := receiver.DBConn.Create(&entity.SFunctionClaim{
		FunctionName: req.FunctionName,
	})

	if result.Error != nil {
		log.Error("FunctionClaimRepository.CreateFunctionClaim: " + result.Error.Error())
		return errors.New("failed to create function claim")
	}

	return nil
}

func (receiver *FunctionClaimRepository) UpdateFunctionClaim(req request.UpdateFunctionClaimRequest) error {
	updateResult := receiver.DBConn.Model(&entity.SFunctionClaim{}).Where("id = ?", req.FunctionClaimID).
		Updates(map[string]interface{}{
			"function_name": req.FunctionName,
		})

	if updateResult.Error != nil {
		log.Error("FunctionClaimRepository.UpdateFunctionClaim: " + updateResult.Error.Error())
		return errors.New("failed to update function claim")
	}

	return nil
}

func (receiver *FunctionClaimRepository) DeleteFunctionClaim(req request.DeleteFunctionClaimRequest) error {
	deleteResult := receiver.DBConn.Delete(&entity.SFunctionClaim{}, req.ID).Error

	if deleteResult != nil {
		log.Error("FunctionClaimRepository.DeleteFunctionClaim: " + deleteResult.Error())
		return errors.New("failed to delete function claim")
	}
	return nil
}
