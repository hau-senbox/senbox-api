package repository

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sen-global-api/internal/domain/entity"
)

type ImageRepository struct {
	DBConn *gorm.DB
}

func NewImageRepository(dbConn *gorm.DB) *ImageRepository {
	return &ImageRepository{DBConn: dbConn}
}

func (receiver *ImageRepository) GetAllByIds(ids []int) ([]entity.SImage, error) {
	var images []entity.SImage
	err := receiver.DBConn.Table(entity.SImage{}.TableName()).Find(&images).Where("id IN (?)", ids).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.SImage{}, nil
		}
		log.Error("ImageRepository.GetAllByIds: " + err.Error())
		return nil, errors.New("failed to get images")
	}

	return images, err
}

func (receiver *ImageRepository) GetAllByName(imageName string) ([]entity.SImage, error) {
	var images []entity.SImage
	err := receiver.DBConn.Table(entity.SImage{}.TableName()).Find(&images).Where("image_name = ", imageName).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.SImage{}, nil
		}
		log.Error("ImageRepository.GetAllByName: " + err.Error())
		return nil, errors.New("failed to get images")
	}

	return images, err
}

func (receiver *ImageRepository) GetByID(id uint64) (*entity.SImage, error) {
	var images entity.SImage
	err := receiver.DBConn.Model(&entity.SImage{}).Where("id = ?", id).First(&images).Error
	if err != nil {
		log.Error("ImageRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get image")
	}
	return &images, nil
}

func (receiver *ImageRepository) GetByKey(key string) (*entity.SImage, error) {
	var images entity.SImage
	err := receiver.DBConn.Model(&entity.SImage{}).Where("`key` = ?", key).First(&images).Error
	if err != nil {
		log.Error("ImageRepository.GetByKey: " + err.Error())
		return nil, errors.New("failed to get image")
	}
	return &images, nil
}

func (receiver *ImageRepository) CreateImages(images []entity.SImage) error {
	db := receiver.DBConn.Begin()
	if err := db.Table(entity.SImage{}.TableName()).Create(&images).Error; err != nil {
		db.Rollback()
		log.Error("ImageRepository.CreateImages: " + err.Error())
		return errors.New("failed to create images")
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		log.Error("ImageRepository.CreateImages: " + err.Error())
		return errors.New("failed to create images")
	}

	return nil
}

func (receiver *ImageRepository) CreateImage(image entity.SImage) error {
	db := receiver.DBConn.Begin()
	if err := db.Table(entity.SImage{}.TableName()).Create(&image).Error; err != nil {
		db.Rollback()
		log.Error("ImageRepository.CreateImages: " + err.Error())
		return errors.New("failed to create images")
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		log.Error("ImageRepository.CreateImages: " + err.Error())
		return errors.New("failed to create images")
	}

	return nil
}

func (receiver *ImageRepository) DeleteImage(id uint64) error {
	db := receiver.DBConn.Begin()

	if err := db.Table(entity.SImage{}.TableName()).Where("id = ?", id).Delete(&entity.SImage{}).Error; err != nil {
		db.Rollback()
		log.Error("ImageRepository.DeleteImage: " + err.Error())
		return errors.New("failed to delete image")
	}

	if err := db.Commit().Error; err != nil {
		db.Rollback()
		log.Error("ImageRepository.DeleteImage: " + err.Error())
		return errors.New("failed to delete image")
	}

	return nil
}
