package repository

import (
	"errors"
	"sen-global-api/internal/domain/entity"

	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type PdfRepository struct {
	DBConn *gorm.DB
}

func NewPdfRepository(dbConn *gorm.DB) *PdfRepository {
	return &PdfRepository{DBConn: dbConn}
}

func (r *PdfRepository) Save(list *entity.SPdf) (error) {
	if err := r.DBConn.Model(&entity.SPdf{}).Create(&list).Error; err != nil {
		return err
	}
	return nil
}

func (receiver *PdfRepository) GetByKey(key string) (*entity.SPdf, error) {
	var pdf entity.SPdf
	err := receiver.DBConn.Model(&entity.SPdf{}).Where("`key` = ?", key).First(&pdf).Error
	if err != nil {
		log.Error("PdfRepository.GetByKey: " + err.Error())
		return nil, errors.New("failed to get image")
	}

	return &pdf, nil
}