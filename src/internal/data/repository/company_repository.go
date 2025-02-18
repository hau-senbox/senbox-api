package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"
	"sen-global-api/internal/domain/request"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type CompanyRepository struct {
	DBConn *gorm.DB
}

func NewCompanyRepository(dbConn *gorm.DB) *CompanyRepository {
	return &CompanyRepository{DBConn: dbConn}
}

func (receiver *CompanyRepository) GetByID(id uint) (*entity.SCompany, error) {
	var company entity.SCompany
	err := receiver.DBConn.Where("id = ?", id).First(&company).Error
	if err != nil {
		log.Error("CompanyRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get company")
	}
	return &company, nil
}

func (receiver *CompanyRepository) CreateCompany(req request.CreateCompanyRequest) error {
	result := receiver.DBConn.Create(&entity.SCompany{
		CompanyName: req.CompanyName,
		Address:     req.Address,
		Description: req.Description,
	})

	if result.Error != nil {
		log.Error("CompanyRepository.CreateCompany: " + result.Error.Error())
		return errors.New("failed to create company")
	}

	return nil
}

func (receiver *CompanyRepository) UpdateCompany(req request.UpdateCompanyRequest) error {
	updateResult := receiver.DBConn.Model(&entity.SCompany{}).Where("id = ?", req.ID).
		Updates(map[string]interface{}{
			"company_name": req.CompanyName,
			"address":      req.Address,
			"description":  req.Description,
		})

	if updateResult.Error != nil {
		log.Error("CompanyRepository.UpdateCompany: " + updateResult.Error.Error())
		return errors.New("failed to update company")
	}

	return nil
}
