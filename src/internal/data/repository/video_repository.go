package repository

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"sen-global-api/internal/domain/entity"
)

type VideoRepository struct {
	DBConn *gorm.DB
}

func NewVideoRepository(dbConn *gorm.DB) *VideoRepository {
	return &VideoRepository{DBConn: dbConn}
}

func (receiver *VideoRepository) GetAllByIDs(ids []int) ([]entity.SVideo, error) {
	var videos []entity.SVideo
	err := receiver.DBConn.Model(entity.SVideo{}).Find(&videos).Where("id IN (?)", ids).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.SVideo{}, nil
		}
		log.Error("VideoRepository.GetAllByIDs: " + err.Error())
		return nil, errors.New("failed to get videos")
	}

	return videos, err
}

func (receiver *VideoRepository) GetAllByName(videoName string) ([]entity.SVideo, error) {
	var videos []entity.SVideo
	err := receiver.DBConn.Model(entity.SVideo{}).Find(&videos).Where("video_name = ", videoName).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return []entity.SVideo{}, nil
		}
		log.Error("VideoRepository.GetAllByName: " + err.Error())
		return nil, errors.New("failed to get videos")
	}

	return videos, err
}

func (receiver *VideoRepository) GetByID(id uint64) (*entity.SVideo, error) {
	var videos entity.SVideo
	err := receiver.DBConn.Model(&entity.SVideo{}).Where("id = ?", id).First(&videos).Error
	if err != nil {
		log.Error("VideoRepository.GetByID: " + err.Error())
		return nil, errors.New("failed to get video")
	}

	return &videos, nil
}

func (receiver *VideoRepository) GetByKey(key string) (*entity.SVideo, error) {
	var videos entity.SVideo
	err := receiver.DBConn.Model(&entity.SVideo{}).Where("`key` = ?", key).First(&videos).Error
	if err != nil {
		log.Error("VideoRepository.GetByKey: " + err.Error())
		return nil, errors.New("failed to get video")
	}

	return &videos, nil
}

func (receiver *VideoRepository) CreateVideos(videos []entity.SVideo) error {
	if err := receiver.DBConn.Model(entity.SVideo{}).Create(&videos).Error; err != nil {
		log.Error("VideoRepository.CreateVideos: " + err.Error())
		return errors.New("failed to create videos")
	}

	return nil
}

func (receiver *VideoRepository) CreateVideo(video entity.SVideo) error {
	if err := receiver.DBConn.Model(entity.SVideo{}).Create(&video).Error; err != nil {
		log.Error("VideoRepository.CreateVideos: " + err.Error())
		return errors.New("failed to create videos")
	}

	return nil
}

func (receiver *VideoRepository) DeleteVideo(id uint64) error {
	if err := receiver.DBConn.Model(entity.SVideo{}).Where("id = ?", id).Delete(&entity.SVideo{}).Error; err != nil {
		log.Error("VideoRepository.DeleteVideo: " + err.Error())
		return errors.New("failed to delete video")
	}

	return nil
}
